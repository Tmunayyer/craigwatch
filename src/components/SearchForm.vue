<style module>
.search-form {
  background-color: #f3f2f2;
  border: 1px solid #a7a7a7;

  box-sizing: border-box;
  width: 100%;
  max-width: 375px;

  display: flex;
  flex-direction: column;

  padding: 1em;
  border-radius: 4px;
}

legend {
  padding: 0.2em;
}

.submit-button {
  width: 100%;
  display: flex;
  justify-content: flex-end;
}
</style>

<template>
  <fieldset class="search-form">
    <legend>new search</legend>
    <InputWithError
      id="name"
      v-model="name"
      :label="'search name:'"
      :placeholder="'appliances'"
      :hasErr="nameErr"
      :errMsg="nameErrMsg"
    />

    <InputWithError
      id="url"
      v-model="url"
      :label="'craigslist search url:'"
      :placeholder="'https://newyork.craigslist.org/d/appliances/search/ppa'"
      :hasErr="urlErr"
      :errMsg="urlErrMsg"
    />
    <div class="submit-button">
      <button v-on:click="handleSubmit">monitor listings</button>
    </div>
  </fieldset>
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
