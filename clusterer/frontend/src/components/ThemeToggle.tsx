import { IconButton, Tooltip } from '@mui/material';
import Brightness4Icon from '@mui/icons-material/Brightness4';
import Brightness7Icon from '@mui/icons-material/Brightness7';

interface ThemeToggleProps {
  isDarkMode: boolean;
  onToggle: () => void;
}

const ThemeToggle = ({ isDarkMode, onToggle }: ThemeToggleProps) => (
  <Tooltip title={isDarkMode ? "Светлая тема" : "Тёмная тема"}>
    <IconButton
      onClick={onToggle}
      color="inherit"
      sx={{
        ml: 2,
        color: 'var(--text-main)',
        backgroundColor: 'transparent',
        '&:hover': {
          backgroundColor: 'var(--hover)',
        }
      }}
    >
      {isDarkMode ? <Brightness7Icon /> : <Brightness4Icon />}
    </IconButton>
  </Tooltip>
);

export default ThemeToggle;
