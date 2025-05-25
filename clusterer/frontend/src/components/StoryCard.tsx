import type { TimelineEvent } from './Timeline';
import {
  Typography, Card, CardContent, CardActionArea,
  Box, Chip, Stack
} from '@mui/material';
import AccessTimeIcon from '@mui/icons-material/AccessTime';
import { format } from 'date-fns';
import { ru } from 'date-fns/locale';

export interface Story { 
  id: number; 
  title: string; 
  summary: string; 
  timeline: TimelineEvent[];
  date?: Date;
  tags?: string[];
}

interface StoryCardProps {
  story: Story;
  onOpen: (story: Story) => void;
  isSelected: boolean;
}

const StoryCard = ({ story, onOpen, isSelected }: StoryCardProps) => (
  <Card 
    sx={{ 
      height: '100%',
      transition: 'all 0.3s ease',
      '&:hover': {
        transform: 'translateY(-2px)',
        boxShadow: '0 4px 8px rgba(0,0,0,0.1)',
      },
      borderRadius: 2,
      overflow: 'hidden',
      border: '1px solid',
      borderColor: isSelected ? 'text.secondary' : 'divider',
    }}
  >
    <CardActionArea 
      onClick={() => onOpen(story)}
      disableRipple
      sx={{ 
        height: '100%',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'stretch',
      }}
    >
      <CardContent sx={{ flexGrow: 1, p: 3 }}>
        <Stack spacing={2}>
          <Box>
            <Typography 
              variant="h6" 
              sx={{
                fontWeight: 600,
                mb: 1,
                lineHeight: 1.3,
                color: 'text.secondary'
              }}
            >
              {story.title}
            </Typography>
            <Typography 
              variant="body2" 
              color="text.secondary"
              sx={{
                display: '-webkit-box',
                WebkitLineClamp: 3,
                WebkitBoxOrient: 'vertical',
                overflow: 'hidden',
                lineHeight: 1.5,
              }}
            >
              {story.summary}
            </Typography>
          </Box>
          
          <Box sx={{ mt: 'auto' }}>
            {story.date && (
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                <AccessTimeIcon sx={{ fontSize: 16, mr: 0.5, color: 'text.secondary' }} />
                <Typography variant="caption" color="text.secondary">
                  {format(story.date, 'd MMMM yyyy', { locale: ru })}
                </Typography>
              </Box>
            )}
            
            {story.tags && story.tags.length > 0 && (
              <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                {story.tags.map((tag, index) => (
                  <Chip
                    key={index}
                    label={tag}
                    size="small"
                    sx={{
                      backgroundColor: 'primary.light',
                      color: 'primary.contrastText',
                      '&:hover': {
                        backgroundColor: 'primary.main',
                      },
                    }}
                  />
                ))}
              </Stack>
            )}
          </Box>
        </Stack>
      </CardContent>
    </CardActionArea>
  </Card>
);

export default StoryCard;