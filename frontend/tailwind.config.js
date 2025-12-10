/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/**/*.{js,ts,jsx,tsx,mdx}"
  ],
  theme: {
    extend: {
      colors: {
        primary: "#333533",
        secondary: '#404040',
        border: '#e0e0e0',
        borderActive: '#c1c1c1'
      },
      borderRadius: {
        primary: '30px',
        input: '12px'
      }
    },
    
  },
  plugins: [],
}

