import { useState } from 'react';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Separator } from '@/components/ui/separator';
import {
  Maximize2, Minimize2, CheckCircle, XCircle, ExternalLink,
  Image as ImageIcon, FileText, Share2, Calendar, Hash
} from 'lucide-react';

export default function ContentPreview({ content, onClose, onStatusUpdate, onDelete, onGenerateImage }) {
  const [isFullScreen, setIsFullScreen] = useState(false);

  if (!content) return null;

  const isBlog = content.type === 'blog';

  return (
    <Dialog open={!!content} onOpenChange={(open) => { if (!open) onClose(); }}>
      <DialogContent
        className={`${isFullScreen ? 'max-w-[95vw] h-[90vh]' : 'max-w-2xl max-h-[85vh]'} rounded-3xl overflow-hidden flex flex-col`}
        data-testid="content-preview-modal"
      >
        <DialogHeader className="pb-3 shrink-0">
          <div className="flex items-center justify-between">
            <div>
              <DialogTitle className="text-lg heading-font">
                {isBlog ? 'Blog Preview' : 'Social Post Preview'}
              </DialogTitle>
              <DialogDescription className="text-xs mt-1">
                Preview your content before publishing
              </DialogDescription>
            </div>
            <Button variant="ghost" size="icon" className="h-8 w-8" onClick={() => setIsFullScreen(!isFullScreen)} data-testid="toggle-fullscreen-btn">
              {isFullScreen ? <Minimize2 className="h-4 w-4" /> : <Maximize2 className="h-4 w-4" />}
            </Button>
          </div>
        </DialogHeader>

        <ScrollArea className="flex-1">
          <div className={isFullScreen ? 'grid grid-cols-1 lg:grid-cols-2 gap-6 p-2' : 'space-y-4 p-2'}>
            {/* Content Details */}
            <div className="space-y-4">
              <div className="flex items-center gap-2 flex-wrap">
                {isBlog
                  ? <Badge className="bg-primary/20 text-primary border-0"><FileText className="h-3 w-3 mr-1" />Blog</Badge>
                  : <Badge className="bg-secondary/20 text-secondary-foreground border-0"><Share2 className="h-3 w-3 mr-1" />Social</Badge>
                }
                <Badge variant="outline" className="capitalize">{content.status || 'draft'}</Badge>
                {content.platform && <Badge variant="outline" className="capitalize">{content.platform}</Badge>}
              </div>

              {content.title && (
                <h2 className="text-xl font-bold heading-font">{content.title}</h2>
              )}

              {content.image_url && (
                <div className="rounded-xl overflow-hidden border">
                  <img src={content.image_url} alt="Content" className="w-full h-48 object-cover" />
                </div>
              )}

              <div
                className="text-sm leading-relaxed text-foreground/90 prose prose-sm max-w-none max-h-96 overflow-y-auto"
                dangerouslySetInnerHTML={{ __html: content.content }}
              />

              {content.hashtags?.length > 0 && (
                <div className="flex items-center gap-1 flex-wrap">
                  <Hash className="h-3 w-3 text-primary" />
                  {content.hashtags.map((tag, i) => (
                    <span key={i} className="text-xs text-primary">{tag}</span>
                  ))}
                </div>
              )}

              {content.created_at && (
                <div className="flex items-center gap-2 text-xs text-muted-foreground">
                  <Calendar className="h-3 w-3" />
                  {new Date(content.created_at).toLocaleDateString()}
                </div>
              )}
            </div>

            {/* Preview Panel (full-screen mode) */}
            {isFullScreen && (
              <div className="border rounded-2xl p-6 bg-muted/20">
                <p className="text-xs uppercase tracking-widest text-muted-foreground mb-4">Live Preview</p>
                {isBlog ? (
                  <div className="space-y-3">
                    <h3 className="text-lg font-bold heading-font">{content.title}</h3>
                    <div
                      className="text-sm leading-relaxed text-muted-foreground prose prose-sm max-w-none max-h-96 overflow-y-auto"
                      dangerouslySetInnerHTML={{ __html: content.content }}
                    />
                  </div>
                ) : (
                  <div className="max-w-sm mx-auto border rounded-xl p-4 bg-card">
                    {content.image_url && <img src={content.image_url} alt="" className="w-full rounded-lg mb-3" />}
                    <div
                      className="text-sm prose prose-sm max-w-none"
                      dangerouslySetInnerHTML={{ __html: content.content }}
                    />
                    {content.hashtags?.length > 0 && <p className="text-xs text-primary mt-2">{content.hashtags.join(' ')}</p>}
                  </div>
                )}
              </div>
            )}
          </div>
        </ScrollArea>

        <Separator className="shrink-0" />

        <div className="flex flex-wrap gap-2 pt-2 shrink-0">
          <Button className="rounded-xl btn-hover" onClick={() => { onStatusUpdate?.(content.id, 'posted'); onClose(); }} data-testid="approve-content-btn">
            <CheckCircle className="h-4 w-4 mr-2" /> Approve & Post
          </Button>
          <Button variant="destructive" className="rounded-xl" onClick={() => { onDelete?.(content.id); onClose(); }} data-testid="reject-content-btn">
            <XCircle className="h-4 w-4 mr-2" /> Reject
          </Button>
          {onGenerateImage && (
            <Button variant="outline" className="rounded-xl" onClick={() => onGenerateImage(content.id)} data-testid="generate-image-btn">
              <ImageIcon className="h-4 w-4 mr-2" /> Generate Image
            </Button>
          )}
          {content.post_url && (
            <Button variant="outline" className="rounded-xl" asChild>
              <a href={content.post_url} target="_blank" rel="noopener noreferrer">
                <ExternalLink className="h-4 w-4 mr-2" /> View Post
              </a>
            </Button>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}
