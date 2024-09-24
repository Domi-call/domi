package settings

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Gpfs struct {
	Api            string `json:"api"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	FileSystemName string `json:"fileSystemName"`
	RootPath       string `json:"rootPath"`
	QuotaLimmit    uint16 `json:"quotaLimmit"`
	QuotaMax       uint16 `json:"quotaMax"`
	Server         string `json:"server"`
	Sysuser        string `json:"sysuser"`
	PrivateKey     string `json:"privateKey"`
}

// 定义结构体来表示 JSON 数据
type Quota struct {
	BlockGrace     string `json:"blockGrace"`
	BlockInDoubt   int    `json:"blockInDoubt"`
	BlockLimit     int    `json:"blockLimit"`
	BlockQuota     int    `json:"blockQuota"`
	BlockUsage     int    `json:"blockUsage"`
	FilesGrace     string `json:"filesGrace"`
	FilesInDoubt   int    `json:"filesInDoubt"`
	FilesLimit     int    `json:"filesLimit"`
	FilesQuota     int    `json:"filesQuota"`
	FilesUsage     int    `json:"filesUsage"`
	FilesetName    string `json:"filesetName"`
	FilesystemName string `json:"filesystemName"`
	IsDefaultQuota bool   `json:"isDefaultQuota"`
	ObjectId       int    `json:"objectId"`
	ObjectName     string `json:"objectName"`
	QuotaId        int    `json:"quotaId"`
	QuotaType      string `json:"quotaType"`
}

type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Response struct {
	Quotas []Quota `json:"quotas"`
	Status Status  `json:"status"`
}

var (
	Api                   = "https://xxx.xxx.xxx.xxx:443"
	FileSystemName        = "gpfsdata"
	RootPath              = "/remote-home/"
	QuotaLimmit    uint16 = 50
	QuotaMax       uint16 = 70
	DefaultServer         = "IP:22"
)

// SendHTTPRequest 发送 HTTP 请求的方法
func SendHTTPRequest(method, url string, requestBody map[string]interface{}, username, password string) ([]byte, error) {
	var jsonData []byte
	var err error

	// 如果请求体不为空，则将请求体转换为 JSON 格式
	if requestBody != nil {
		jsonData, err = json.Marshal(requestBody)
		if err != nil {
			return nil, fmt.Errorf("error marshalling JSON: %v", err)
		}
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(username, password)

	// 跳过证书验证（不推荐用于生产环境）
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// 发送 HTTP 请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return responseBody, nil
}

// 在GPFS中创建用户目录, 必须先创建用户, 然后才能将fileset分给用户
func CreateFileset(username string, g Gpfs) error {
	// API URL
	url := fmt.Sprintf("%s/scalemgmt/v2/filesystems/%s/filesets", g.Api, g.FileSystemName)
	// 请求体
	requestBody := map[string]interface{}{
		"filesetName": username,
		"path":        g.RootPath + username,
		"permissions": "700",
		"owner":       username + ":" + username,
	}

	// 发送 HTTP 请求
	res, err := SendHTTPRequest("POST", url, requestBody, g.Username, g.Password)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	fmt.Println("Response:", string(res))

	resMap := make(map[string]interface{})
	err = json.Unmarshal(res, &resMap)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	statusStr := resMap["status"]
	statusMap := statusStr.(map[string]interface{})
	code := statusMap["code"].(float64)
	codeStr := fmt.Sprintf("%v", code)
	if strings.HasPrefix(codeStr, "20") {
		return nil
	} else {
		return fmt.Errorf("code=%v,message: %s", statusMap["code"], statusMap["message"])
	}
}

// 设置用户home dir
func SetUserHomeDir(username string, g Gpfs) error {
	session, err := GetSessionGpfsRemote(g)
	if err != nil {
		return err
	}
	defer session.Close()
	// Command to set the home directory
	homeDir := g.RootPath + username
	setHomeDirCmd := fmt.Sprintf("usermod -d %s %s", homeDir, username)
	_, err = RunCommandAtRemoteServer(setHomeDirCmd, g)
	if err != nil {
		log.Printf("Error setting home directory: %s", err)
		return err
	}

	cp := fmt.Sprintf("cp -r /etc/skel/. %s", homeDir)
	_, err = RunCommandAtRemoteServer(cp, g)
	if err != nil {
		log.Printf("Error copying files: %" + err.Error())
	}
	//chownSnapshots := fmt.Sprintf("mv %s/.snapshots %s/%s.snapshots", homeDir, g.RootPath, username)
	//if err = runCommand(chownSnapshots); err != nil {
	//	log.Printf("Error moving .snapshots: %s", err)
	//}
	//再把用户目录的权限设置为用户
	//chownHome := fmt.Sprintf("chown -R %s:%s %s", username, username, homeDir)
	chownHome := fmt.Sprintf("chown -R %s:%s %s", username, username, homeDir)
	_, err = RunCommandAtRemoteServer(chownHome, g)
	if err != nil {
		log.Printf("Error chown home directory: %s", err)
	}
	//chownSnapshotsRoot := fmt.Sprintf("mv %s/%s.snapshots %s/.snapshots", g.RootPath, username, homeDir)
	//if err = runCommand(chownSnapshotsRoot); err != nil {
	//	log.Printf("Error moving .snapshots: %s", err)
	//}

	return nil
}

// 创建或修改目录配额
func SetQuota(filesetName string, blockSoftLimit uint16, blockHardLimit uint16, g Gpfs) error {
	// API URL
	url := fmt.Sprintf("%s/scalemgmt/v2/filesystems/%s/quotas", g.Api, g.FileSystemName)

	// 请求体
	requestBody := map[string]interface{}{
		"operationType":  "setQuota",
		"quotaType":      "FILESET",
		"objectName":     filesetName,
		"blockSoftLimit": fmt.Sprintf("%dG", blockSoftLimit),
		"blockHardLimit": fmt.Sprintf("%dG", blockHardLimit),
	}

	// 发送 HTTP 请求
	res, err := SendHTTPRequest("POST", url, requestBody, g.Username, g.Password)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	fmt.Println("Response:", string(res))

	resMap := make(map[string]interface{})
	err = json.Unmarshal(res, &resMap)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	statusStr := resMap["status"]
	statusMap := statusStr.(map[string]interface{})
	code := statusMap["code"].(float64)
	codeStr := fmt.Sprintf("%v", code)
	if strings.HasPrefix(codeStr, "20") {
		return nil
	} else {
		return fmt.Errorf("code=%v,message: %s", statusMap["code"], statusMap["message"])
	}
}

// 查看目录限额
func QueryFilesetQuota(filesetName string, g Gpfs) (error, []Quota) {
	api := fmt.Sprintf("%s/scalemgmt/v2/filesystems/%s/quotas", g.Api, g.FileSystemName)

	// 如果查询条件filesetName不为空, 则查询filesetName
	if filesetName != "" {
		filter := fmt.Sprintf("objectName=%s", filesetName)
		encodedFilter := url.QueryEscape(filter)
		api = fmt.Sprintf("%s/scalemgmt/v2/filesystems/%s/quotas?filter=%s", g.Api, g.FileSystemName, encodedFilter)
	}
	// API URL
	responseBody, err := SendHTTPRequest("GET", api, nil, g.Username, g.Password)
	if err != nil {
		fmt.Println("Error:", err)
		return err, nil
	}
	fmt.Println("Response:", string(responseBody))
	var response Response
	err = json.Unmarshal([]byte(responseBody), &response)
	if err != nil {
		fmt.Printf(" QueryUserFilesetQuota Parsed error")
		return err, nil
	}
	return nil, response.Quotas
}

// 查询文件集使用情况
func QueryFilesetUsage(filesetName string, g Gpfs) (error, int64) {
	path := g.RootPath + filesetName
	cmd := fmt.Sprintf("du -s %s | awk '{print $1}'", path)
	session, err := GetSessionGpfsRemote(g)
	if err != nil {
		fmt.Println("Error:", err)
		return err, 0
	}
	defer session.Close()
	output, err := RunCommandAtRemoteServer(cmd, g)
	if err != nil {
		fmt.Println("Error:", err)
		return err, 0
	}
	strs := strings.Split(strings.TrimSpace(string(output)), "\n")
	//仅返回是数字的一行
	str := strs[len(strs)-1]
	float, err := strconv.ParseFloat(str, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return err, 0
	}
	return nil, int64(float) * 1024
}

func RunCommandAtRemoteServer(cmd string, g Gpfs) ([]byte, error) {
	session, err := GetSessionGpfsRemote(g)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer session.Close()
	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return nil, fmt.Errorf("command failed: %s, output: %s, error: %v", cmd, output, err)
	}
	return output, nil
}

func GetSessionGpfsRemote(g Gpfs) (*ssh.Session, error) {
	signer, err := ssh.ParsePrivateKey([]byte(g.PrivateKey))
	if err != nil {
		return nil, err
	}
	// SSH client configuration
	config := &ssh.ClientConfig{
		User: g.Sysuser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	// Connect to the SSH server
	client, err := ssh.Dial("tcp", g.Server, config)
	if err != nil {
		return nil, err
	}
	session, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}
	return session, nil
}
