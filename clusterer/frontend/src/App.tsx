import React, { useState } from 'react';
import {
  Typography, Box, AppBar, Toolbar, Tab, Tabs, IconButton, useTheme, useMediaQuery,
  Container, Fade, Grow
} from '@mui/material';
import StoryCard from './components/StoryCard';
import Timeline from './components/Timeline';
import type { Story } from './components/StoryCard';
import { STORIES } from './Data'; // Предполагается, что STORIES содержит поле 'id'

function App() {
  const [selectedStoryId, setSelectedStoryId] = useState<number | null>(null);
  const [tabValue, setTabValue] = useState(0);
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));

  const selectedStory = STORIES.find(story => story.id === selectedStoryId);

  return (
    <Box sx={{
      minHeight: '100vh',
      width: '100%',
      bgcolor: 'background.default',
      background: 'linear-gradient(180deg, #f5f7fa 0%, #ffffff 100%)',
      display: 'flex',
      flexDirection: 'column'
    }}>
      <AppBar
        position="sticky"
        elevation={0}
        sx={{
          backdropFilter: 'blur(8px)',
          backgroundColor: 'rgba(255, 255, 255, 0.9)',
          borderBottom: '1px solid rgba(0, 0, 0, 0.12)'
        }}
      >
        <Toolbar>
          <Tabs
            value={tabValue}
            onChange={(_, v) => setTabValue(v)}
            textColor="primary"
            indicatorColor="primary"
            sx={{
              '& .MuiTab-root': {
                fontSize: '1rem',
                fontWeight: 600,
                textTransform: 'none',
                minWidth: 120,
                transition: 'all 0.3s ease',
                '&:hover': {
                  color: 'primary.main',
                }
              }
            }}
          >
            <Tab label="Читать" />
            <Tab label="Избранное" disabled />
            <Tab label="Настройки" disabled />
          </Tabs>
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
        boxSizing: 'border-box',
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
                background: 'linear-gradient(45deg, #2196F3 30%, #21CBF3 90%)',
                backgroundClip: 'text',
                textFillColor: 'transparent',
                WebkitBackgroundClip: 'text',
                WebkitTextFillColor: 'transparent'
              }}
            >
              Сюжеты дня
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
                  <Typography variant="h5" sx={{ fontWeight: 600, mb: 2 }}>
                    {selectedStory.title}
                  </Typography>
                  <Typography variant="body1" color="text.secondary" sx={{ mb: 3, lineHeight: 1.7 }}>
                    {selectedStory.summary}
                  </Typography>
                  <Timeline events={selectedStory.timeline} />
                </Box>
              </Fade>
            ) : (
              <Typography variant="h6" align="center" color="text.secondary" sx={{ mt: 5 }}>
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
             <Typography variant="h6" sx={{ fontWeight: 600, mb: 2 }}>
                {selectedStory.title}
              </Typography>
              <Typography variant="body1" color="text.secondary" sx={{ mb: 3, lineHeight: 1.7 }}>
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
