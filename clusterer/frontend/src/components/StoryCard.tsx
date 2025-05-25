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
  onOpen: (id: number) => void;
  isSelected: boolean;
}

// Теги
const StoryTags = ({ tags }: { tags: string[] }) => (
  <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
    {tags.map(tag => (
      <Chip
        key={tag}
        label={tag}
        size="small"
        sx={{
          backgroundColor: 'var(--primary-light)',
          color: 'var(--primary-contrast)',
          fontWeight: 500,
          '&:hover': { backgroundColor: 'var(--primary)' },
        }}
      />
    ))}
  </Stack>
);

// Дата
const StoryDate = ({ date }: { date: Date }) => (
  <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
    <AccessTimeIcon sx={{ fontSize: 16, mr: 0.5, color: 'var(--text-secondary)' }} />
    <Typography variant="caption" sx={{ color: 'var(--text-secondary)' }}>
      {format(date, 'd MMMM yyyy', { locale: ru })}
    </Typography>
  </Box>
);

const StoryCard = ({ story, onOpen, isSelected }: StoryCardProps) => (
  <Card
    sx={{
      height: '100%',
      border: `1px solid ${isSelected ? 'var(--primary)' : 'var(--border)'}`,
      backgroundColor: 'var(--bg)',
      color: 'var(--text-main)',
      boxShadow: 'none',
      borderRadius: 2, // Без скруглений!
      transition: 'border-color 0.3s',
      position: 'relative',
      cursor: 'pointer',
      '&:hover': {
        borderColor: 'var(--primary)',
        boxShadow: 'none',
      }
    }}
    onClick={() => onOpen(story.id)}
    tabIndex={0}
    role="button"
  >
    <CardActionArea
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
                color: 'var(--text-main)'
              }}
            >
              {story.title}
            </Typography>
            <Typography
              variant="body2"
              sx={{
                display: '-webkit-box',
                WebkitLineClamp: 3,
                WebkitBoxOrient: 'vertical',
                overflow: 'hidden',
                lineHeight: 1.5,
                color: 'var(--text-secondary)'
              }}
            >
              {story.summary}
            </Typography>
          </Box>
          <Box sx={{ mt: 'auto' }}>
            {story.date && <StoryDate date={story.date} />}
            {story.tags && story.tags.length > 0 && <StoryTags tags={story.tags} />}
          </Box>
        </Stack>
      </CardContent>
    </CardActionArea>
  </Card>
);

export default StoryCard;
