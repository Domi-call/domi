package settings

import (
	"fmt"
	"log"
	"os/exec"

	"golang.org/x/crypto/ssh"
)

// checkUserExists 检查用户是否存在
func checkUserExists(username string) bool {
	cmd := exec.Command("getent", "passwd", username)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// unlockUser 解锁用户账户
func unlockUser(username string) error {
	cmd := exec.Command("passwd", "-u", username)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to unlock user: %w", err)
	}
	return nil
}

// setPassword 设置用户密码
func setPassword(username, password string) error {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("echo '%s:%s' | chpasswd", username, password))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set user password: %w", err)
	}
	return nil
}

// 通过SSH Key创建用户 sysuser必须有 SUDO 权限
func (s *Settings) CreateUser(username, password, serverRoot string, g Gpfs) error {
	signer, err := ssh.ParsePrivateKey([]byte(g.PrivateKey))
	if err != nil {
		return err
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
		return err
	}
	defer client.Close()

	// Helper function to run a command on the remote server
	runCommand := func(cmd string) error {
		session, err := client.NewSession()
		if err != nil {
			return fmt.Errorf("failed to create session: %w", err)
		}
		defer session.Close()

		output, err := session.CombinedOutput(cmd)
		if err != nil {
			return fmt.Errorf("command failed: %s, output: %s, error: %w", cmd, output, err)
		}
		return nil
	}

	// Command to create a new user
	createUserCmd := fmt.Sprintf("useradd %s", username)
	if err := runCommand(createUserCmd); err != nil {
		log.Printf("Error creating user: %s", err)
		return err
	}

	// Command to set the password
	setPasswordCmd := fmt.Sprintf("echo '%s:%s' | chpasswd", username, password)
	if err := runCommand(setPasswordCmd); err != nil {
		log.Printf("Error setting password: %s", err)
		return err
	}

	// Command to unlock the user
	unlockUserCmd := fmt.Sprintf("usermod -U %s", username)
	if err := runCommand(unlockUserCmd); err != nil {
		log.Printf("Error unlocking user: %s", err)
		return err
	}

	return nil
}
