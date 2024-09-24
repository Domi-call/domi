import { fetchURL, fetchJSON } from "./utils";

// 获取Quota
export function getQuota(param: string) {
  return fetchJSON<Quota[]>(`/api/gpfs/getquota?objectName=${param}`, {
    method: "GET",
  });
}

// 获取默认Quota
export function getQuotaDefault() {
  return fetchJSON<QuotaDefault>(`/api/gpfs/getQuotaDefault`, {
    method: "GET",
  });
}

// 获取用户文件集使用量
export function getUserUsage() {
  return fetchJSON<Usage>(`/api/gpfs/getUserUsage`, {
    method: "GET",
  });
}

// 获取用户文件集限额
export function getUserQuota() {
  return fetchJSON<QuotaDefault>(`/api/gpfs/getUserQuota`, {
    method: "GET",
  });
}

// 保存Quota
export async function setQuota(params: object) {
  await fetchURL(`/api/gpfs/setquota`, {
    method: "POST",
    body: JSON.stringify(params),
  });
}

// 创建Fileset
export async function createFileset() {
  await fetchURL(`/api/gpfs/fileset`, {
    method: "POST",
  });
}
