<template>
  <div class="card floating">
    <div class="card-title">
      <h2>{{ $t("prompts.quotaSetting") }}</h2>
    </div>

    <div class="card-content">
      <p>
        <strong>{{ $t("prompts.displayName") }}</strong> {{ name }}
      </p>
      <p>
        <strong class="form-title">{{ $t("prompts.quotaSoft") }}:</strong>
        <input class="input" type="number" v-model="soft" /> GB
      </p>
      <p>
        <strong class="form-title">{{ $t("prompts.quotaHard") }}:</strong>
        <input class="input" type="number" v-model="hard" /> GB
      </p>
    </div>

    <div class="card-action">
      <button
        id="focus-prompt"
        @click="closeHovers"
        class="button button--flat"
        :aria-label="$t('buttons.cancel')"
        :title="$t('buttons.cancel')"
      >
        {{ $t("buttons.cancel") }}
      </button>
      <button
        id="focus-prompt"
        type="submit"
        @click="saveQuota"
        class="button button--flat"
        :aria-label="$t('buttons.save')"
        :title="$t('buttons.save')"
      >
        {{ $t("buttons.save") }}
      </button>
    </div>
  </div>
</template>

<script>
import { gpfs } from "@/api";
import { mapActions, mapState } from "pinia";
import { useFileStore } from "@/stores/file";
import { useLayoutStore } from "@/stores/layout";

export default {
  name: "quota",
  inject: ["$showError", "$showSuccess"],
  data() {
    return {
      soft: "",
      hard: "",
    };
  },

  computed: {
    ...mapState(useFileStore, [
      "req",
      "selected",
      "selectedCount",
      "isListing",
    ]),
    name: function () {
      return this.selectedCount === 0
        ? this.req.name
        : this.req.items[this.selected[0]].name;
    },
  },
  methods: {
    ...mapActions(useLayoutStore, ["closeHovers"]),
    getQuota: async function () {
      const quotas = await gpfs.getQuota(this.name);
      if (quotas.length > 0) {
        this.soft = quotas[0].blockQuota / 1024 / 1024;
        this.hard = quotas[0].blockLimit / 1024 / 1024;
      }
    },
    saveQuota: async function () {
      this.closeHovers();
      // 比较 soft 和 hard, soft 必须小于 hard
      if (this.soft > this.hard) {
        this.$showError("Quota : soft mast less than hard");
        return;
      }
      const params = {
        filesetName: this.name,
        quotaLimmit: this.soft,
        quotaMax: this.hard,
      };
      await gpfs.setQuota(params);
      this.$showSuccess("保存成功");
    },
  },
  mounted() {
    // TODO 增加加载动画
    this.getQuota();
  },
};
</script>
<style scoped>
.form-title {
  margin-right: 10px;
}
</style>
