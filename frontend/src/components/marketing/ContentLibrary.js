import { useState, useEffect, useCallback } from 'react';
import api, { addUserIdToParams } from '@/lib/api';
import { toast } from 'sonner';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Skeleton } from '@/components/ui/skeleton';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Eye, Trash2, RefreshCw, FileText, Share2, ExternalLink } from 'lucide-react';
import ContentPreview from '@/components/marketing/ContentPreview';

const STATUS_COLORS = {
  draft: 'bg-muted text-muted-foreground',
  generated: 'bg-secondary/20 text-secondary-foreground',
  posted: 'bg-primary/20 text-primary',
  failed: 'bg-destructive/20 text-destructive',
};

export default function ContentLibrary() {
  const [content, setContent] = useState([]);
  const [loading, setLoading] = useState(true);
  const [typeFilter, setTypeFilter] = useState('all');
  const [statusFilter, setStatusFilter] = useState('all');
  const [previewContent, setPreviewContent] = useState(null);

  const fetchContent = useCallback(async () => {
    setLoading(true);
    try {
      const params = addUserIdToParams({});
      if (typeFilter !== 'all') params.type = typeFilter;
      if (statusFilter !== 'all') params.status = statusFilter;
      const response = await api.get('/content', { params });
      setContent(response.data?.data || []);
    } catch {
      // silently handle
    } finally {
      setLoading(false);
    }
  }, [typeFilter, statusFilter]);

  useEffect(() => { fetchContent(); }, [fetchContent]);

  const handleDelete = async (id) => {
    try {
      await api.delete(`/content/${id}`, { params: addUserIdToParams() });
      toast.success('Content deleted');
      fetchContent();
    } catch {
      toast.error('Failed to delete');
    }
  };

  const handleStatusUpdate = async (id, status) => {
    try {
      await api.put(`/content/${id}/status`, { status }, { params: addUserIdToParams() });
      toast.success(`Status updated to ${status}`);
      fetchContent();
    } catch {
      toast.error('Failed to update');
    }
  };

  return (
    <>
      <Card className="rounded-2xl" data-testid="content-library">
        <CardHeader className="pb-3 border-b">
          <div className="flex items-center justify-between">
            <CardTitle className="text-lg heading-font">Content Library</CardTitle>
            <Button variant="ghost" size="icon" onClick={fetchContent} className="h-8 w-8" data-testid="refresh-content-btn">
              <RefreshCw className="h-4 w-4" />
            </Button>
          </div>
          <div className="flex gap-2 mt-2">
            <Select value={typeFilter} onValueChange={setTypeFilter}>
              <SelectTrigger className="w-24 h-8 text-xs rounded-lg" data-testid="type-filter"><SelectValue /></SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Types</SelectItem>
                <SelectItem value="blog">Blog</SelectItem>
                <SelectItem value="social">Social</SelectItem>
              </SelectContent>
            </Select>
            <Select value={statusFilter} onValueChange={setStatusFilter}>
              <SelectTrigger className="w-28 h-8 text-xs rounded-lg" data-testid="status-filter"><SelectValue /></SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Status</SelectItem>
                <SelectItem value="draft">Draft</SelectItem>
                <SelectItem value="generated">Generated</SelectItem>
                <SelectItem value="posted">Posted</SelectItem>
                <SelectItem value="failed">Failed</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </CardHeader>
        <ScrollArea className="h-[480px]">
          <CardContent className="p-3 space-y-2">
            {loading ? (
              Array.from({ length: 3 }).map((_, i) => <Skeleton key={i} className="h-24 rounded-xl" />)
            ) : content.length === 0 ? (
              <div className="text-center py-12 text-muted-foreground">
                <FileText className="h-8 w-8 mx-auto mb-2 opacity-50" />
                <p className="text-sm">No content yet. Generate some using the chat!</p>
              </div>
            ) : (
              content.map(item => (
                <div key={item.id} className="p-3 rounded-xl border bg-card hover:shadow-sm transition-shadow">
                  <div className="flex items-start justify-between gap-2">
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2 mb-1">
                        {item.type === 'blog' ? <FileText className="h-3 w-3 text-primary shrink-0" /> : <Share2 className="h-3 w-3 text-secondary shrink-0" />}
                        <span className="text-sm font-medium truncate">{item.title || item.content?.substring(0, 40) || 'Untitled'}</span>
                      </div>
                      <p className="text-xs text-muted-foreground line-clamp-2">{item.content?.substring(0, 100)}</p>
                      <div className="flex items-center gap-2 mt-2">
                        <Badge className={`text-[10px] ${STATUS_COLORS[item.status] || ''}`}>{item.status}</Badge>
                        {item.platform && <Badge variant="outline" className="text-[10px] capitalize">{item.platform}</Badge>}
                      </div>
                    </div>
                  </div>
                  <div className="flex gap-1 mt-2">
                    <Button size="sm" variant="ghost" className="h-7 text-xs" onClick={() => setPreviewContent(item)} data-testid={`lib-preview-${item.id}`}>
                      <Eye className="h-3 w-3 mr-1" /> View
                    </Button>
                    <Button size="sm" variant="ghost" className="h-7 text-xs" onClick={() => handleStatusUpdate(item.id, 'posted')}>
                      <ExternalLink className="h-3 w-3 mr-1" /> Post
                    </Button>
                    <Button size="sm" variant="ghost" className="h-7 text-xs text-destructive hover:text-destructive" onClick={() => handleDelete(item.id)}>
                      <Trash2 className="h-3 w-3 mr-1" /> Delete
                    </Button>
                  </div>
                </div>
              ))
            )}
          </CardContent>
        </ScrollArea>
      </Card>

      {previewContent && (
        <ContentPreview content={previewContent} onClose={() => setPreviewContent(null)} onStatusUpdate={handleStatusUpdate} onDelete={handleDelete} />
      )}
    </>
  );
}
