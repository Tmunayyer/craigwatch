<style scoped>
/* Tooltip container */
.tooltip {
  position: relative;
  display: inline-block;
  text-decoration: underline;
  text-decoration-style: dotted;
  font-weight: bold;
}

/* Tooltip text */
.tooltip .tooltiptext {
  /* visibility: hidden; */
  /* width: 120px; */
  background-color: #555;
  color: #fff;
  text-align: center;
  padding: 5px 0;
  border-radius: 6px;

  /* Position the tooltip text */
  position: absolute;
  z-index: 1;
  bottom: 125%;
  left: 50%;
  margin-left: -60px;

  /* Fade in tooltip */
  opacity: 1;
  transition: opacity 0.3s;
}

/* Tooltip arrow */
.tooltip .tooltiptext::after {
  content: "";
  position: absolute;
  top: -36%;
  left: 15%;
  margin-left: -5px;
  border-width: 5px;

  border-style: solid;
  border-color: #555 transparent transparent transparent;

  transform: rotate(180deg);
}

/* Show the tooltip text when you mouse over the tooltip container */
.tooltip:hover .tooltiptext {
  visibility: visible;
  opacity: 1;
}
</style>

<template>
  <div class="tooltip">
    (?)
    <span class="tooltiptext">{{text}}</span>
  </div>
</template>

<script>
export default {
  name: "Tooltip",
  props: ["text"]
};
</script>