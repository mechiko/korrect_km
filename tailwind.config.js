/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./alcogo/public/assets/css/daisyui.full.css", "./webapp/**/*.{css,html}"],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
  ]
}
