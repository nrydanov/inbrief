import { useState, useEffect } from 'react';
import {
  Typography, Box, AppBar, Toolbar, Tab, Tabs, useTheme, useMediaQuery, Fade, Grow
} from '@mui/material';
import StoryCard from './components/StoryCard';
import Timeline from './components/Timeline';
import ThemeToggle from './components/ThemeToggle';
import { STORIES } from './Data'; // Предполагается, что STORIES содержит поле 'id'

function App() {
  const [selectedStoryId, setSelectedStoryId] = useState<number | null>(null);
  const [tabValue, setTabValue] = useState(0);
  const [isDarkMode, setIsDarkMode] = useState(() => {
    const saved = localStorage.getItem('darkMode');
    return saved ? JSON.parse(saved) : false;
  });
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));

  const selectedStory = STORIES.find(story => story.id === selectedStoryId);

  useEffect(() => {
    document.documentElement.setAttribute('data-theme', isDarkMode ? 'dark' : 'light');
    localStorage.setItem('darkMode', JSON.stringify(isDarkMode));
  }, [isDarkMode]);

  return (
    <Box sx={{
      minHeight: '100vh',
      width: '100%',
      maxWidth: '100vw',
      overflowX: 'hidden',
      bgcolor: 'background.default',
      background: 'var(--gradient-background)',
      display: 'flex',
      flexDirection: 'column'
    }}>
      <AppBar
        position="sticky"
        elevation={0}
        sx={{
          backdropFilter: 'blur(8px)',
          backgroundColor: isDarkMode ? 'rgba(18, 18, 18, 0.9)' : 'rgba(255, 255, 255, 0.9)',
          borderBottom: '1px solid var(--color-border)'
        }}
      >
        <Toolbar sx={{ justifyContent: 'space-between' }}>
          <Tabs
            value={tabValue}
            onChange={(_, v) => setTabValue(v)}
            textColor="primary"
            indicatorColor="primary"
            sx={{
              '& .MuiTabs-root': {
                color: 'var(--color-tab-text)',
              },
              '& .MuiTab-root': {
                fontSize: '1rem',
                fontWeight: 600,
                textTransform: 'none',
                minWidth: 120,
                transition: 'all 0.3s ease',
                color: 'var(--color-tab-text)',
                '&:hover': {
                  color: 'var(--color-primary)',
                },
                '&.Mui-selected': {
                  color: 'var(--color-primary)',
                },
                '&.Mui-disabled': {
                  color: 'var(--color-text-disabled)',
                }
              }
            }}
          >
            <Tab label="Читать" />
            <Tab label="Избранное" disabled />
            <Tab label="Настройки" disabled />
          </Tabs>
          <ThemeToggle isDarkMode={isDarkMode} onToggle={() => setIsDarkMode(!isDarkMode)} />
        </Toolbar>
      </AppBar>
      <Box sx={{
        flexGrow: 1,
        py: { xs: 2, md: 4 },
        px: { xs: 2, md: 4 },
        display: 'flex',
        flexDirection: isMobile ? 'column' : 'row',
        gap: 3,
        width: '100%',
        maxWidth: '100%',
        boxSizing: 'border-box',
        overflowX: 'hidden'
      }}>
        {/* Левая часть: Список сюжетов */}
        <Box sx={{
          flexBasis: isMobile ? 'auto' : '50%',
          flexShrink: 0,
          overflowY: 'auto',
          pr: isMobile ? 0 : 3, // Увеличиваем отступ справа только на десктопе
          pb: isMobile ? 3 : 0, // Увеличиваем отступ снизу только на мобайле
        }}>
           <Fade in timeout={800}>
            <Typography
              variant={isMobile ? "h5" : "h4"}
              align="center"
              sx={{
                fontWeight: 800,
                mb: 4,
                color: 'var(--color-text-primary)'
              }}
            >
              Сюжеты
            </Typography>
          </Fade>
          <Box sx={{
             display: 'grid',
             gap: 3,
             // На десктопе одна колонка для списка, на мобайле - одна колонка
             gridTemplateColumns: '1fr'
          }}>
            {STORIES.map((story, index) => (
              <Grow in timeout={500} style={{ transitionDelay: `${index * 100}ms` }} key={story.id}>
                <Box>
                  <StoryCard
                    story={story}
                    // При клике просто устанавливаем выбранный ID сюжета
                    onOpen={() => setSelectedStoryId(story.id)}
                    isSelected={story.id === selectedStoryId}
                  />
                </Box>
              </Grow>
            ))}
          </Box>
        </Box>

        {/* Правая часть: Таймлайн выбранного сюжета */}
        {!isMobile && (
          <Box sx={{
            flexBasis: '50%',
            flexShrink: 0,
            overflowY: 'auto',
            pl: 3, // Увеличиваем отступ слева
          }}>
            {selectedStory ? (
              <Fade in timeout={500}>
                <Box>
                  <Typography variant="h5" sx={{ fontWeight: 600, mb: 2, color: 'var(--color-text-primary)' }}>
                    {selectedStory.title}
                  </Typography>
                  <Typography variant="body1" sx={{ mb: 3, lineHeight: 1.7, color: 'var(--color-text-secondary)' }}>
                    {selectedStory.summary}
                  </Typography>
                  <Timeline events={selectedStory.timeline} />
                </Box>
              </Fade>
            ) : (
              <Typography variant="h6" align="center" sx={{ mt: 5, color: 'var(--color-text-secondary)' }}>
                Выберите сюжет, чтобы увидеть таймлайн
              </Typography>
            )}
          </Box>
        )}

        {/* На мобайле таймлайн пока не показываем в этом макете, так как нет двух колонок */}
        {isMobile && selectedStory && (
           <Box sx={{
            flexGrow: 1,
             overflowY: 'auto',
             mt: 3,
             pt: 3,
             borderTop: '1px solid rgba(0, 0, 0, 0.12)'
           }}>
             <Typography variant="h6" sx={{ fontWeight: 600, mb: 2, color: 'var(--color-text-primary)' }}>
                {selectedStory.title}
              </Typography>
              <Typography variant="body1" sx={{ mb: 3, lineHeight: 1.7, color: 'var(--color-text-secondary)' }}>
                {selectedStory.summary}
              </Typography>
             <Timeline events={selectedStory.timeline} />
           </Box>
        )}

      </Box>
    </Box>
  );
}

export default App;
