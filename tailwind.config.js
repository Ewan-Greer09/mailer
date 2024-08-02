/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./services/frontend/views/**/*.templ"],
  theme: {
    extend: {},
  },
  plugins: [require("@tailwindcss/typography"), require("daisyui")],
  daisyui: {
    themes: [
      {
        mytheme: {
          primary: "#6b21a8",
          secondary: "#a78bfa",
          accent: "#e879f9",
          neutral: "#1c1917",
          "base-100": "#111827",
          info: "#d946ef",
          success: "#00ff00",
          warning: "#f59e0b",
          error: "#ff0000",
        },
      },
    ],
  },
};
