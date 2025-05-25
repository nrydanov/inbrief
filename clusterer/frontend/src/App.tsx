import React, { useState } from 'react';
import {
  Typography, Dialog, DialogTitle,
  DialogContent, Box, AppBar, Toolbar, Tab, Tabs, IconButton, useTheme, useMediaQuery
} from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';
import StoryCard from './components/StoryCard';
import Timeline from './components/Timeline';
import type { Story } from './components/StoryCard';
import { STORIES } from './Data';

function App() {
  const [open, setOpen] = useState(false);
  const [selectedStory, setSelectedStory] = useState<Story | null>(null);
  const [tabValue, setTabValue] = useState(0);
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));

  return (
    <Box sx={{ width: '100vw', height: '100vh', bgcolor: 'background.default' }}>
      <AppBar position="static">
        <Toolbar>
          <Tabs value={tabValue} onChange={(_, v) => setTabValue(v)} textColor="inherit" indicatorColor="secondary">
            <Tab label="ЧИТАТЬ" />
            <Tab label="ИЗБРАННОЕ" disabled />
            <Tab label="НАСТРОЙКИ" disabled />
          </Tabs>
        </Toolbar>
      </AppBar>
      <Box sx={{ maxWidth: 1000, mx: 'auto', px: 3, py: 3 }}>
        <Typography variant={isMobile ? "h5" : "h4"} align="center" sx={{ fontWeight: 'bold', mb: 3 }}>
          Сюжеты дня
        </Typography>
        {STORIES.map((story) => (
          <StoryCard key={story.id} story={story} onOpen={(s) => { setSelectedStory(s); setOpen(true); }} />
        ))}
        <Dialog
          open={open}
          onClose={() => { setOpen(false); setSelectedStory(null); }}
          fullWidth
          maxWidth="xl"
          fullScreen={isMobile}
          PaperProps={{ sx: { borderRadius: isMobile ? 0 : 2, maxHeight: '90vh', height: '90vh' } }}
        >
          <DialogTitle sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Typography variant="h6">{selectedStory?.title}</Typography>
            <IconButton onClick={() => { setOpen(false); setSelectedStory(null); }} size="small">
              <CloseIcon />
            </IconButton>
          </DialogTitle>
          <DialogContent>
            <Typography sx={{ mb: 2 }}>{selectedStory?.summary}</Typography>
            {selectedStory && <Timeline events={selectedStory.timeline} />}
          </DialogContent>
        </Dialog>
      </Box>
    </Box>
  );
}

export default App;
