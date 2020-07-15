<style scoped>
.search-field {
  margin-bottom: 0.5em;
}

.search-label {
  width: 100%;
  font-weight: 600;
}

.search-input {
  box-sizing: border-box;
  width: 100%;
}

.error {
  padding: 2px 2px 2px 2px;
  font-size: 0.8em;
  color: red;
}
</style>

<template>
  <div class="search-field">
    <div class="search-label">{{ label }}</div>
    <input
      class="search-input"
      :placeholder="placeholder"
      type="text"
      v-bind:value="value"
      v-on:input="$emit('input', $event.target.value)"
    />
    <div class="error" v-if="hasErr">{{ errMsg }}</div>
  </div>
</template>

<script>
export default {
  name: "InputWithError",
  props: ["value", "label", "placeholder", "hasErr", "errMsg"]
};
</script>