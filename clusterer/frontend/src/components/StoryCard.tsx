import type { TimelineEvent } from './Timeline';
import {
  Typography, Card, CardContent, CardActionArea,
  Box, Chip, Stack
} from '@mui/material';
import AccessTimeIcon from '@mui/icons-material/AccessTime';
import { format } from 'date-fns';
import { ru } from 'date-fns/locale';

// Story interface representing a story object
export interface Story {
  id: number;
  title: string;
  summary: string;
  timeline: TimelineEvent[];
  date?: Date;
  tags?: string[];
}

// Props interface for StoryCard component
interface StoryCardProps {
  story: Story;
  onOpen: (id: number) => void;
  isSelected: boolean;
}

// Style objects for reuse and easier maintenance
const cardSx = {
  height: '100%',
  transition: 'all 0.3s ease',
  '&:hover': {
    transform: 'translateY(-2px)',
    boxShadow: '0 4px 8px rgba(0,0,0,0.1)',
  },
  borderRadius: 2,
  overflow: 'hidden',
  border: '1px solid',
};

const titleSx = {
  fontWeight: 600,
  mb: 1,
  lineHeight: 1.3,
  color: 'text.secondary',
};

const summarySx = {
  display: '-webkit-box',
  WebkitLineClamp: 3,
  WebkitBoxOrient: 'vertical',
  overflow: 'hidden',
  lineHeight: 1.5,
  color: 'text.secondary',
};

// Subcomponent for rendering story tags
const StoryTags = ({ tags }: { tags: string[] }) => (
  <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
    {tags.map(tag => (
      <Chip key={tag} label={tag} size="small" sx={{
        backgroundColor: 'primary.light',
        color: 'primary.contrastText',
        '&:hover': { backgroundColor: 'primary.main' },
      }} />
    ))}
  </Stack>
);

// Subcomponent for rendering story date
const StoryDate = ({ date }: { date: Date }) => (
  <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
    <AccessTimeIcon sx={{ fontSize: 16, mr: 0.5, color: 'text.secondary' }} />
    <Typography variant="caption" color="text.secondary">
      {format(date, 'd MMMM yyyy', { locale: ru })}
    </Typography>
  </Box>
);

// Main StoryCard component
const StoryCard = ({ story, onOpen, isSelected }: StoryCardProps) => (
  <Card sx={{ ...cardSx, borderColor: isSelected ? 'text.secondary' : 'divider' }}>
    <CardActionArea
      onClick={() => onOpen(story.id)}
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
            <Typography variant="h6" sx={titleSx}>{story.title}</Typography>
            <Typography variant="body2" sx={summarySx}>{story.summary}</Typography>
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
