import { createVuetify } from "vuetify";

import "@mdi/font/css/materialdesignicons.css";
import "vuetify/styles";

export default defineNuxtPlugin((app) => {
  const vuetify = createVuetify({
    theme: {
      defaultTheme: "dark",
    },
  });
  app.vueApp.use(vuetify);
});
