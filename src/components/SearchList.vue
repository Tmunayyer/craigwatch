<template>
  <h1>search list</h1>
</template>

<script>
export default {
  name: "SearchList",
  data() {
    return {
      name: "",
      url: ""
    };
  },
  methods: {
    redirector: function(search) {
      this.$router.push(`/result/${search.ID}`);
    },
    getSearchList: async function() {
      const apiUrl = `/api/v1/search`;
      const response = await fetch(apiUrl);
      const body = await response.json();

      console.log("the body:", body);
    }
  }
};
</script>