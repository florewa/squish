/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{svelte,js,ts}', './index.html'],
  theme: {
    extend: {
      colors: {
        surface: {
          50:  '#f5f5f5',
          100: '#ebebeb',
          800: '#1e1e1e',
          900: '#141414',
          950: '#0d0d0d',
        },
        accent: {
          400: '#a78bfa',
          500: '#8b5cf6',
          600: '#7c3aed',
        }
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
      }
    },
  },
  plugins: [],
}
