<script setup>
import SearchForm from './SearchForm.vue';
import {onMounted, ref, watch} from 'vue';
import _ from 'lodash';
import DateRangePicker from './DateRangePicker.vue';
import moment from 'moment';

const searchFilter = ref('');
const dateRangeFilter = ref({ startDate: '', endDate: '' });
const currentPage = ref(1);
const pageSize = ref(5); 
const totalData = ref(0)
const totalAmount = ref(0)
const totalPage = ref(0)

const data = ref({
    success: false,
    data: {
        total_data: 0,
        total_amount: '',
        total_page: 0,
        current_page: 0,
        next_page: false,
        data: []
    },
    error: {
        code: "",
        message: ""
    }
});



const filteredData = ref([]);

onMounted(async() => {
    const response = await fetch("api/order");
    data.value = await response.json();

    filteredData.value = data.value.data.data;
    currentPage.value = data.value.data.current_page;
    totalData.value = data.value.data.total_data;
    totalAmount.value = data.value.data.total_amount;
    totalPage.value = data.value.data.total_page;
})


const fetchData = async () => {
  let url = `http://localhost:8080/order`;

  const params = new URLSearchParams();
  if (searchFilter.value) {
    params.append('keyword', searchFilter.value);
  }
  if (dateRangeFilter.value.startDate && dateRangeFilter.value.endDate) {
    params.append('date_start', dateRangeFilter.value.startDate);
    params.append('date_end', dateRangeFilter.value.endDate);
  }

  params.append('limit', pageSize.value)
  params.append('page', currentPage.value)

  console.log(dateRangeFilter)
  try {
    const response = await fetch(`${url}?${params.toString()}`);
    const resp = await response.json();
    filteredData.value = resp.data.data;
    totalData.value = resp.data.total_data;
    totalAmount.value = resp.data.total_amount;
    totalPage.value = resp.data.total_page;
  } catch (error) {
    console.error('Error fetching data:', error);
    filteredData.value = [];
  }
};

const debouncedSearch = _.debounce((search) => {
  searchFilter.value = search;
  fetchData();
}, 300);

const handleSearch = (search) => {
  debouncedSearch(search);
};


const debouncedDateRange = _.debounce((dateRange) => {
    dateRangeFilter.value = dateRange;
  fetchData();
}, 300);

const handleDateRange = (dateRange) => {
    debouncedDateRange(dateRange);
};

watch(() => data.data, (newData) => {
    if (searchFilter.value === '' && !dateRangeFilter.value.start) {
    filteredData.value = newData;
  }
});

const nextPage = () => {
  if (currentPage.value < totalPage.value) {
    currentPage.value++;
    fetchData();
  }
};

const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--;
    fetchData();
  }
};
</script>

<template>
    <div class="flex flex-col items-center container mx-auto w-full"> 
        <div class="bg-white relative border rounded-lg">
            <div class="flex w-full">
                <SearchForm @search="handleSearch"/>
            </div>

            <div class="w-80 mb-8">
                <DateRangePicker @dateRange="handleDateRange" />
            </div>

            <div class="w-80 mb-5 px-4 text-sm text-gray-700">
            Total Amount: <b>{{ totalAmount }}</b>
            </div>

            <table class="w-full text-sm text-left text-gray-500">
                <thead class="text-xs text-gray-700">
                    <tr>
                        <th class="px-4 py-3">Order Name</th>
                        <th class="px-4 py-3">Customer Company</th>
                        <th class="px-4 py-3">Customer Name</th>
                        <th class="px-4 py-3">Order Date</th>
                        <th class="px-4 py-3">Delivered Amount</th>
                        <th class="px-4 py-3">Total Amount</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="singleData in filteredData" :key="singleData?.order_name">
                        <td class="px-4 py-3 font-medium">{{ singleData?.order_name }}</td>
                        <td class="px-4 py-3 font-medium">{{ singleData?.customer_company }}</td>
                        <td class="px-4 py-3 font-medium">{{ singleData?.customer_name }}</td>
                        <td class="px-4 py-3 font-medium">{{ moment(singleData?.order_date).format('MMM Do YY hh:mm A') }}</td>
                        <td class="px-4 py-3 font-medium">{{ singleData?.delivered_amount }}</td>
                        <td class="px-4 py-3 font-medium">{{ singleData?.total_amount }}</td>
                    </tr>
                </tbody>

            </table>
        </div>
        <div class="mt-5">
        <button class="px-4 py-3" @click="prevPage" :disabled="currentPage.value === 1">
            <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 48 48"><path fill="none" stroke="#000" stroke-linecap="round" stroke-linejoin="round" stroke-width="4" d="M31 36L19 24L31 12"/></svg>
        </button>
        <span class="px-4 py-2 border-2 text-sm	 rounded border-violet-500">{{ currentPage }}</span>
        <button class="px-4 py-3" @click="nextPage" :disabled="currentPage.value === totalPage">
            <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 48 48"><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="4" d="m19 12l12 12l-12 12"/></svg>
        </button>
        </div>
    </div>
</template>
