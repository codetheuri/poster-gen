/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}", // This line tells Tailwind to scan all your Vue components
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}