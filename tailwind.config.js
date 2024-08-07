/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./services/frontend/views/**/*.templ"],
  theme: {
    extend: {},
  },
  plugins: [require("@tailwindcss/typography"), require("daisyui")],
  daisyui: {
    themes: ["night", "dark", "sunset", "black"],
  },
};
