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
            color: '#475569', // slate-600
            h1: { color: '#1e293b' }, // slate-900
            h2: { color: '#1e293b' },
            h3: { color: '#1e293b' },
            strong: { color: '#2563eb' }, // blue-600
            code: { color: '#2563eb' },
            blockquote: {
              color: '#64748b',
              borderLeftColor: '#3b82f6',
              backgroundColor: '#f8fafc',
              padding: '0.5rem 1rem',
              borderRadius: '0.5rem',
            },
            table: {
              fontSize: '0.75rem',
              lineHeight: '1rem',
              borderCollapse: 'separate',
              borderSpacing: '0',
              width: '100%',
              marginTop: '1.5rem',
              marginBottom: '1.5rem',
              borderRadius: '0.75rem',
              overflow: 'hidden',
              border: '1px solid #e2e8f0',
            },
            thead: {
              backgroundColor: '#f1f5f9',
              borderBottomWidth: '1px',
              borderBottomColor: '#e2e8f0',
            },
            'thead th': {
              color: '#475569',
              fontWeight: '700',
              padding: '0.75rem',
              textAlign: 'left',
            },
            'tbody tr': {
              backgroundColor: '#ffffff',
              borderBottomWidth: '1px',
              borderBottomColor: '#f1f5f9',
            },
            'tbody tr:nth-child(even)': {
              backgroundColor: '#f8fafc',
            },
            'tbody td': {
              padding: '0.75rem',
              color: '#334155',
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
