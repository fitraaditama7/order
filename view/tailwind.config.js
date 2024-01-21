/** @type {import('tailwindcss').Config} */
import colors from "tailwindcss/colors";
import formsPlugin from '@tailwindcss/forms';

export default {
  purge: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}', './node_modules/vue-tailwind-datepicker/**/*.js'],
  darkMode: false, // or 'media' or 'class'
  theme: {
    extend: {
      colors: {
        "vtd-primary": colors.sky,
        "vtd-secondary": colors.white,
      },
    },
  },
  variants: {
    extend: {},
  },
  plugins: [
    formsPlugin()
  ],
}

