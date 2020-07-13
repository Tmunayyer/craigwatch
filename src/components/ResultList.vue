<style module>
.result-list {
  padding: 0.2em 0.2em 0.2em 0.2em;

  box-sizing: border-box;
  width: 100%;
  max-width: 625px;
}

.result-listitem {
  box-shadow: 0 0px 2px 0 rgba(0, 0, 0, 0.3);
  transition: 0.3s;
  padding: 0.5em;
  border-radius: 4px;
  margin-bottom: 0.5em;
}

.listitem-title {
  font-size: 1em;
  font-weight: bold;
}

.listitem-date {
  font-size: 0.8em;
}

.listitem-body {
  padding: 0.2em 0.2em 0.2em 0.2em;
  overflow: hidden;
}

.price {
  font-weight: bold;
}
</style>

<template>
  <ul class="result-list">
    <li v-for="(listing) in resultList" v-bind:key="listing.URL">
      <div class="result-listitem">
        <div class="listitem-header">
          <div class="listitem-title">{{ listing.Title }}</div>
          <div class="listitem-date">{{ formatDate(listing.UnixDate*1000) }}</div>
        </div>
        <hr />
        <div class="listitem-body">
          <div class="price">${{ listing.Price }}</div>
          <br />
          <a class="result-header-url" v-bind:href="listing.Link" target="_blank">{{ listing.Link }}</a>
        </div>
      </div>
    </li>
  </ul>
</template>

<script>
export default {
  name: "ResultList",
  props: ["resultList"],
  methods: {
    formatDate: function(unixTimestamp) {
      var options = {
        year: "numeric",
        month: "numeric",
        day: "numeric",
        hour: "2-digit",
        minute: "2-digit"
      };
      var today = new Date(unixTimestamp);

      return today.toLocaleDateString("en-US", options);
    }
  }
};
</script>