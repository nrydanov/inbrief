import type { TimelineEvent } from './Timeline';

export interface Story { 
    id: number; 
    title: string; 
    summary: string; 
    timeline: TimelineEvent[];
}

import {
    Typography, Card, CardContent, CardActionArea
  } from '@mui/material';

const StoryCard = ({ story, onOpen }: { story: Story; onOpen: (story: Story) => void }) => (
    <Card sx={{ mb: 2, borderColor: 'divider', borderRadius: 1 }}>
      <CardActionArea onClick={() => onOpen(story)}>
        <CardContent>
          <Typography variant="h6">{story.title}</Typography>
        </CardContent>
      </CardActionArea>
    </Card>
  );


export default StoryCard;