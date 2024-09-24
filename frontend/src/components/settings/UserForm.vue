<template>
  <div>
    <p v-if="!isDefault && props.user !== null">
      <label for="username">{{ t("settings.username") }}</label>
      <input
        class="input input--block"
        type="text"
        v-model="user.username"
        id="username"
      />
    </p>

    <p v-if="!isDefault">
      <label for="password">{{ t("settings.password") }}</label>
      <input
        class="input input--block"
        type="password"
        :placeholder="passwordPlaceholder"
        v-model="user.password"
        id="password"
      />
    </p>

    <p>
      <label for="scope">{{ t("settings.scope") }}</label>
      <input
        :disabled="createUserDirData ?? false"
        :placeholder="scopePlaceholder"
        class="input input--block"
        type="text"
        v-model="user.scope"
        id="scope"
      />
    </p>

    <p>
      <label>{{ $t("settings.quotaSoft") }}</label>
      <input
        :disabled="createUserDirData ?? false"
        :placeholder="String(quotaDefault.soft)"
        class="input"
        type="number"
        v-model="user.quotaSoft"
        id="quotaSoft"
      />
      GB
    </p>
    <p>
      <label>{{ $t("settings.quotaHard") }}</label>
      <input
        :disabled="createUserDirData ?? false"
        :placeholder="String(quotaDefault.hard)"
        class="input"
        type="number"
        v-model="user.quotaHard"
        id="quotaHard"
      />
      GB
    </p>

    <p class="small" v-if="displayHomeDirectoryCheckbox">
      <input type="checkbox" v-model="createUserDirData" />
      {{ t("settings.createUserHomeDirectory") }}
    </p>

    <p>
      <label for="locale">{{ t("settings.language") }}</label>
      <languages
        class="input input--block"
        id="locale"
        v-model:locale="user.locale"
      ></languages>
    </p>

    <p v-if="!isDefault && user.perm">
      <input
        type="checkbox"
        :disabled="user.perm.admin"
        v-model="user.lockPassword"
      />
      {{ t("settings.lockPassword") }}
    </p>

    <permissions v-model:perm="user.perm" />
    <commands v-if="enableExec" v-model:commands="user.commands" />

    <div v-if="!isDefault">
      <h3>{{ t("settings.rules") }}</h3>
      <p class="small">{{ t("settings.rulesHelp") }}</p>
      <rules v-model:rules="user.rules" />
    </div>
  </div>
</template>

<script setup lang="ts">
import Languages from "./Languages.vue";
import Rules from "./Rules.vue";
import Permissions from "./Permissions.vue";
import Commands from "./Commands.vue";
import { enableExec } from "@/utils/constants";
import { computed, onMounted, ref, watch } from "vue";
import { useI18n } from "vue-i18n";
import { gpfs } from "@/api/index.js";

const { t } = useI18n();

const createUserDirData = ref<boolean | null>(null);
const originalUserScope = ref<string | null>(null);

const props = defineProps<{
  user: IUserForm;
  isNew: boolean;
  isDefault: boolean;
  createUserDir?: boolean;
}>();

const quotaDefault = ref({
  soft: 0,
  hard: 0,
});

onMounted(() => {
  fetchQuotaDefaultValues();
  if (props.user.scope) {
    originalUserScope.value = props.user.scope;
    createUserDirData.value = props.createUserDir;
  }
});

const fetchQuotaDefaultValues = async () => {
  try {
    const quota = await gpfs.getQuotaDefault();
    quotaDefault.value.soft = quota.soft;
    quotaDefault.value.hard = quota.hard;
    props.user.quotaSoft = props.user.quotaSoft ?? quota.soft;
    props.user.quotaHard = props.user.quotaHard ?? quota.hard;
  } catch (error) {
    console.error("Error fetching default values:", error);
  }
};

const passwordPlaceholder = computed(() =>
  props.isNew ? "" : t("settings.avoidChanges")
);
const scopePlaceholder = computed(() =>
  createUserDirData.value ? t("settings.userScopeGenerationPlaceholder") : ""
);
const displayHomeDirectoryCheckbox = computed(
  () => props.isNew && createUserDirData.value
);

watch(
  () => props.user,
  () => {
    if (!props.user?.perm?.admin) return;
    props.user.lockPassword = false;
  }
);

watch(createUserDirData, () => {
  if (props.user?.scope) {
    props.user.scope = createUserDirData.value
      ? ""
      : originalUserScope.value ?? "";
  }
});
</script>
