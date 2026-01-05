/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        'stock-red': '#f5222d',
        'stock-green': '#52c41a',
      },
      typography: {
        DEFAULT: {
          css: {
            maxWidth: 'none',
            color: '#94a3b8',
            h1: { color: '#f1f5f9' },
            h2: { color: '#f1f5f9' },
            h3: { color: '#f1f5f9' },
            strong: { color: '#f1f5f9' },
            code: { color: '#3b82f6' },
            blockquote: {
              color: '#94a3b8',
              borderLeftColor: '#3b82f6',
            },
            table: {
              fontSize: '0.75rem',
              lineHeight: '1rem',
              borderCollapse: 'collapse',
              width: '100%',
              marginTop: '1.5rem',
              marginBottom: '1.5rem',
            },
            thead: {
              backgroundColor: '#1e293b',
              borderBottomWidth: '1px',
              borderBottomColor: '#334155',
            },
            'thead th': {
              color: '#f1f5f9',
              fontWeight: '600',
              padding: '0.5rem 0.75rem',
              textAlign: 'left',
            },
            'tbody tr': {
              borderBottomWidth: '1px',
              borderBottomColor: '#1e293b',
            },
            'tbody tr:nth-child(even)': {
              backgroundColor: '#0f172a',
            },
            'tbody td': {
              padding: '0.5rem 0.75rem',
              color: '#cbd5e1',
            },
          },
        },
      },
    },
  },
  plugins: [
    require('@tailwindcss/typography'),
  ],
}
