<script setup>
import { ref, watch } from "vue";
import VueTailwindDatepicker from "vue-tailwind-datepicker";

const dateValue = ref({
  startDate: "",
  endDate: "",
});
const formatter = ref({
  date: "YYYY-MM-DD",
  month: "MMM",
});

const emit = defineEmits(['dateRange']);
watch(dateValue, (newValue) => {
    emit('dateRange', newValue);
});
</script>

<template>
    <div class="px-4 flex items-center text-sm font-semibold text-gray-700">
        Created date
    </div>
    <div class="px-4 flex items-center">
        <vue-tailwind-datepicker :formatter="formatter" v-model="dateValue" placeholder="Start data - End date"/>
    </div>
</template>