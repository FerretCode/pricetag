/** @type {import('tailwindcss').Config} */

const colors = require('tailwindcss/colors')

module.exports = {
  content: ["./ui/web/**/*.tmpl"],
  theme: {
    colors: {
      black: colors.black,
      white: colors.white,
      stone: colors.stone,
      blue: colors.blue,
      red: colors.red,
      green: colors.green,
      yellow: colors.yellow,
    },
  },
  plugins: [],
}

