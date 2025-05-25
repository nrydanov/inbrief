import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { ThemeProvider, createTheme } from '@mui/material/styles';
import './theme.css'
import './index.css'
import App from './App.tsx'
import '@fontsource-variable/inter';

const theme = createTheme({
  typography: {
    fontFamily: ['"Inter Variable"', 'sans-serif'].join(','),
  },
});

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ThemeProvider theme={theme}>
      <App />
    </ThemeProvider>
  </StrictMode>,
)