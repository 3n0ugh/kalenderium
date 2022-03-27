<template>
  <v-app>
    <CalendarComponent />
  </v-app>
</template>

<script>
import CalendarComponent from "../components/CalendarComponent";

export default {
  name: "CalendarView",

  components: {
    CalendarComponent,
  },
};
</script>
