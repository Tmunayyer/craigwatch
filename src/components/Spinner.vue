<style scoped>
#spinner-container {
  height: 100%;
}
</style>

<template>
  <div id="spinner-container">
    <ClipLoader :loading="sharedState.loading" />
  </div>
</template>

<script>
import ClipLoader from "vue-spinner/src/ClipLoader.vue";

export const spinnerState = {
  state: {
    loading: false,
    timeout: null,
  },
  setLoading(isLoading) {
    if (isLoading) {
      this.state.timeout = setTimeout(() => {
        if (this.state.timeout !== null) {
          this.state.loading = true;
        }
      }, 500);
    } else {
      this.state.loading = isLoading;
      clearTimeout(this.state.timeout);
      this.state.timeout = null;
    }
  },
};

export default {
  name: "Spinner",
  components: {
    ClipLoader,
  },
  data() {
    return {
      sharedState: spinnerState.state,
    };
  },
};
</script>