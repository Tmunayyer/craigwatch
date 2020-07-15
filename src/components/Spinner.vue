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
    loading: false
  },
  setLoading(isLoading) {
    this.state.loading = isLoading;
  }
};

export default {
  name: "Spinner",
  components: {
    ClipLoader
  },
  data() {
    return {
      sharedState: spinnerState.state
    };
  }
};
</script>