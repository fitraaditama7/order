import { createApp } from 'vue';
import App from './App.vue';
import './style.css';
import VueTailwindDatepicker from 'vue-tailwind-datepicker'


const app = createApp(App);
app.use(VueTailwindDatepicker)

app.mount('#app');