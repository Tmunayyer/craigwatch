<style module>
.search-form {
  box-sizing: border-box;
  max-width: 250px;

  display: flex;
  flex-direction: column;
}

.submit-button {
  width: fit-content;
  align-self: flex-end;
}
</style>

<template>
  <div class="search-form">
    <InputWithError
      id="name"
      v-model="name"
      :label="'search name:'"
      :hasErr="nameErr"
      :errMsg="nameErrMsg"
    />
    <InputWithError
      id="url"
      v-model="url"
      :label="'craigslist search url:'"
      :hasErr="urlErr"
      :errMsg="urlErrMsg"
    />
    <button class="submit-button" v-on:click="handleSubmit">monitor listings</button>
  </div>
</template>

<script>
import InputWithError from "./InputWithError.vue";

export default {
  name: "SearchForm",
  components: {
    InputWithError
  },
  data() {
    return {
      name: "",
      nameErr: false,
      nameErrMsg: "",
      url: "",
      urlErr: false,
      urlErrMsg: ""
    };
  },
  methods: {
    redirector: function(search) {
      this.$router.push(`/result/${search.ID}`);
    },
    handleSubmit: async function() {
      const isInvalid = this.validation();
      if (isInvalid) {
        return;
      }

      const apiUrl = `/api/v1/search`;
      const apiOptions = {
        method: "POST",
        body: JSON.stringify({ Name: this.name, URL: this.url })
      };

      const search = await this.$http(apiUrl, apiOptions);

      this.redirector(search);
    },
    validation: function() {
      const errorMessage = {
        nameRequired: "A name is required.",
        urlRequired: "A url is required.",
        invalidURL: "Invalid craigslist url provided."
      };

      let hasError = false;

      if (this.name === "") {
        hasError = true;
        this.nameErr = true;
        this.nameErrMsg = errorMessage.nameRequired;
      }

      if (this.url === "") {
        hasError = true;
        this.urlErr = true;
        this.urlErrMsg = errorMessage.urlRequired;
      }

      const urlRegex = new RegExp(/https:\/\/.*\.craigslist.org/);
      if (!urlRegex.test(this.url)) {
        hasError = true;
        this.urlErr = true;
        this.urlErrMsg = errorMessage.invalidURL;
      }

      return hasError;
    }
  }
};
</script>
