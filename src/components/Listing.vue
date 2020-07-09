<template>
  <div id="listingItem">
    <h2 style="text-decoration: bold">{{ listing.Title }}</h2>
    <h4>posted on: {{ listing.UnixDate }}</h4>
    <p>{{ listing.Link }}</p>
    <h4>price: ${{ listing.Price }}</h4>
  </div>
</template>

<script>
export default {
  name: "Listing",
  props: ["listing", "index"]
};
</script>
