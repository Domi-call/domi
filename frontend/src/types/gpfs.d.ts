interface Quota {
  blockGrace: string;
  blockInDoubt: number;
  blockLimit: number;
  blockQuota: number;
  blockUsage: number;
  filesGrace: string;
  filesInDoubt: number;
  filesLimit: number;
  filesQuota: number;
  filesUsage: number;
  filesetName: string;
  filesystemName: string;
  isDefaultQuota: boolean;
  objectId: number;
  objectName: string;
  quotaId: number;
  quotaType: string;
}

interface QuotaParams {
  filesetName: string;
  quotaLimmit: number;
  quotaMax: number;
}

interface QuotaDefault {
  soft: number;
  hard: number;
}

interface Usage {
  usage: number;
}
