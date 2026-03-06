import { useState, useRef, useEffect } from 'react';
import api, { addUserIdToParams } from '@/lib/api';
import { toast } from 'sonner';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Checkbox } from '@/components/ui/checkbox';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Send, Loader2, FileText, Share2, Eye, Trash2, Image as ImageIcon, ExternalLink, Bot } from 'lucide-react';
import ContentPreview from '@/components/marketing/ContentPreview';

const PLATFORMS = ['instagram', 'linkedin', 'twitter', 'facebook'];

export default function ContentChat() {
  const [messages, setMessages] = useState([]);
  const [prompt, setPrompt] = useState('');
  const [contentType, setContentType] = useState('blog');
  const [platforms, setPlatforms] = useState([]);
  const [tone, setTone] = useState('');
  const [loading, setLoading] = useState(false);
  const [previewContent, setPreviewContent] = useState(null);
  const bottomRef = useRef(null);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const togglePlatform = (p) => {
    setPlatforms(prev => prev.includes(p) ? prev.filter(x => x !== p) : [...prev, p]);
  };

  const handleGenerate = async () => {
    if (!prompt.trim()) return;
    if (contentType === 'social' && platforms.length === 0) {
      toast.error('Select at least one platform for social posts');
      return;
    }

    const userMsg = { role: 'user', content: prompt, type: contentType, platforms: [...platforms], tone };
    setMessages(prev => [...prev, userMsg]);
    setLoading(true);
    const currentPrompt = prompt;
    setPrompt('');

    try {
      let response;
      if (contentType === 'blog') {
        response = await api.post('/content/blog', {
          topic: currentPrompt,
          ...(tone && { tone }),
          length: 'medium',
        }, { params: addUserIdToParams() });
        setMessages(prev => [...prev, { role: 'assistant', data: response.data.data, type: 'blog' }]);
      } else {
        response = await api.post('/content/social', {
          topic: currentPrompt,
          platforms,
          ...(tone && { tone }),
        }, { params: addUserIdToParams() });
        const data = Array.isArray(response.data.data) ? response.data.data : [response.data.data];
        setMessages(prev => [...prev, { role: 'assistant', data, type: 'social' }]);
      }
      toast.success('Content generated!');
    } catch (error) {
      toast.error(error.response?.data?.message || 'Failed to generate content');
      setMessages(prev => [...prev, { role: 'assistant', type: 'error', error: error.response?.data?.message || 'Generation failed' }]);
    } finally {
      setLoading(false);
    }
  };

  const handleStatusUpdate = async (contentId, status) => {
    try {
      await api.put(`/content/${contentId}/status`, { status }, { params: addUserIdToParams() });
      toast.success(`Content ${status}!`);
    } catch {
      toast.error('Failed to update status');
    }
  };

  const handleDelete = async (contentId) => {
    try {
      await api.delete(`/content/${contentId}`, { params: addUserIdToParams() });
      toast.success('Content deleted');
    } catch {
      toast.error('Failed to delete');
    }
  };

  const handleGenerateImage = async (contentId) => {
    try {
      toast.info('Generating image... This may take up to 1 minute.');
      await api.post(`/content/${contentId}/generate-image`, {
        image_prompt: 'Professional, high-quality image relevant to this content',
      }, { params: addUserIdToParams() });

      toast.success('Image generated successfully!');

      // Reload the content to show the new image
      const updatedContent = await api.get(`/content/${contentId}`, { params: addUserIdToParams() });

      // Update preview if it's open
      if (previewContent && previewContent.id === contentId) {
        setPreviewContent(updatedContent.data.data);
      }

      // Update the message in the chat to show the image
      setMessages(prev => prev.map(msg => {
        if (msg.data?.id === contentId) {
          return { ...msg, data: { ...msg.data, image_url: updatedContent.data.data.image_url } };
        }
        return msg;
      }));
    } catch (error) {
      toast.error('Failed to generate image: ' + (error.response?.data?.error || 'Unknown error'));
    }
  };

  return (
    <>
      <Card className="rounded-2xl flex flex-col h-[600px]" data-testid="content-chat">
        <CardHeader className="pb-3 border-b shrink-0">
          <CardTitle className="text-lg flex items-center gap-2 heading-font">
            <Bot className="h-5 w-5 text-primary" />
            Content Generator
          </CardTitle>
        </CardHeader>

        {/* Messages */}
        <div className="flex-1 overflow-auto p-4 custom-scrollbar">
          <div className="space-y-4">
            {messages.length === 0 && (
              <div className="text-center text-muted-foreground py-12">
                <Bot className="h-12 w-12 mx-auto mb-3 opacity-30" />
                <p className="text-sm">Start by entering a topic to generate content</p>
                <p className="text-xs mt-1 text-muted-foreground/70">Choose Blog or Social, then type your topic</p>
              </div>
            )}

            {messages.map((msg, i) => (
              <div key={i} className={`chat-message flex ${msg.role === 'user' ? 'justify-end' : 'justify-start'}`}>
                {msg.role === 'user' ? (
                  <div className="max-w-[80%] bg-primary text-primary-foreground rounded-2xl rounded-tr-md px-4 py-3">
                    <div className="flex items-center gap-2 mb-1">
                      <Badge variant="outline" className="text-[10px] border-primary-foreground/30 text-primary-foreground">
                        {msg.type === 'blog' ? 'Blog' : 'Social'}
                      </Badge>
                      {msg.platforms?.map(p => (
                        <Badge key={p} variant="outline" className="text-[10px] border-primary-foreground/30 text-primary-foreground capitalize">{p}</Badge>
                      ))}
                    </div>
                    <p className="text-sm">{msg.content}</p>
                  </div>
                ) : msg.type === 'error' ? (
                  <div className="max-w-[80%] bg-destructive/10 border border-destructive/20 rounded-2xl rounded-tl-md px-4 py-3">
                    <p className="text-sm text-destructive">{msg.error}</p>
                  </div>
                ) : msg.type === 'blog' ? (
                  <div className="max-w-[85%] bg-card border rounded-2xl rounded-tl-md px-4 py-3 space-y-2">
                    <p className="text-sm font-semibold">{msg.data?.title || 'Blog Post'}</p>
                    <p className="text-xs text-muted-foreground line-clamp-3">{msg.data?.content?.substring(0, 200)}...</p>
                    <div className="flex flex-wrap gap-2 pt-1">
                      <Button size="sm" variant="outline" className="h-7 text-xs rounded-lg" onClick={() => setPreviewContent(msg.data)} data-testid={`preview-btn-${i}`}>
                        <Eye className="h-3 w-3 mr-1" /> Preview
                      </Button>
                      <Button size="sm" variant="outline" className="h-7 text-xs rounded-lg" onClick={() => handleStatusUpdate(msg.data?.id, 'posted')} data-testid={`approve-btn-${i}`}>
                        <ExternalLink className="h-3 w-3 mr-1" /> Approve
                      </Button>
                      <Button size="sm" variant="outline" className="h-7 text-xs rounded-lg text-destructive hover:text-destructive" onClick={() => handleDelete(msg.data?.id)} data-testid={`delete-btn-${i}`}>
                        <Trash2 className="h-3 w-3 mr-1" /> Reject
                      </Button>
                      <Button size="sm" variant="outline" className="h-7 text-xs rounded-lg" onClick={() => handleGenerateImage(msg.data?.id)}>
                        <ImageIcon className="h-3 w-3 mr-1" /> Image
                      </Button>
                    </div>
                  </div>
                ) : (
                  <div className="max-w-[85%] space-y-2">
                    {(Array.isArray(msg.data) ? msg.data : [msg.data]).map((post, j) => (
                      <div key={j} className="bg-card border rounded-2xl rounded-tl-md px-4 py-3 space-y-2">
                        <div className="flex items-center gap-2">
                          <Badge className="capitalize text-[10px]">{post?.platform || 'social'}</Badge>
                          <Badge variant="outline" className="text-[10px]">{post?.status || 'generated'}</Badge>
                        </div>
                        <p className="text-xs text-muted-foreground line-clamp-3">{post?.content?.substring(0, 150)}</p>
                        {post?.hashtags?.length > 0 && (
                          <p className="text-xs text-primary">{post.hashtags.join(' ')}</p>
                        )}
                        <div className="flex flex-wrap gap-2 pt-1">
                          <Button size="sm" variant="outline" className="h-7 text-xs rounded-lg" onClick={() => setPreviewContent(post)} data-testid={`preview-social-${i}-${j}`}>
                            <Eye className="h-3 w-3 mr-1" /> Preview
                          </Button>
                          <Button size="sm" variant="outline" className="h-7 text-xs rounded-lg" onClick={() => handleStatusUpdate(post?.id, 'posted')}>
                            <ExternalLink className="h-3 w-3 mr-1" /> Approve
                          </Button>
                          <Button size="sm" variant="outline" className="h-7 text-xs rounded-lg text-destructive hover:text-destructive" onClick={() => handleDelete(post?.id)}>
                            <Trash2 className="h-3 w-3 mr-1" /> Reject
                          </Button>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            ))}

            {loading && (
              <div className="flex justify-start">
                <div className="bg-card border rounded-2xl rounded-tl-md px-4 py-3 flex items-center gap-2">
                  <Loader2 className="h-4 w-4 animate-spin text-primary" />
                  <span className="text-sm text-muted-foreground">Generating content...</span>
                </div>
              </div>
            )}
            <div ref={bottomRef} />
          </div>
        </div>

        {/* Input Area */}
        <div className="border-t p-4 space-y-3 shrink-0">
          <div className="flex flex-wrap items-center gap-2">
            <Select value={contentType} onValueChange={setContentType}>
              <SelectTrigger className="w-28 h-8 text-xs rounded-lg" data-testid="content-type-select">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="blog"><FileText className="h-3 w-3 inline mr-1" />Blog</SelectItem>
                <SelectItem value="social"><Share2 className="h-3 w-3 inline mr-1" />Social</SelectItem>
              </SelectContent>
            </Select>

            {contentType === 'social' && (
              <div className="flex items-center gap-3">
                {PLATFORMS.map(p => (
                  <label key={p} className="flex items-center gap-1.5 text-xs cursor-pointer">
                    <Checkbox checked={platforms.includes(p)} onCheckedChange={() => togglePlatform(p)} data-testid={`platform-${p}`} />
                    <span className="capitalize">{p}</span>
                  </label>
                ))}
              </div>
            )}

            <Input value={tone} onChange={e => setTone(e.target.value)} placeholder="Tone (optional)" className="w-32 h-8 text-xs rounded-lg" data-testid="tone-input" />
          </div>

          <div className="flex gap-2">
            <Input
              value={prompt}
              onChange={e => setPrompt(e.target.value)}
              placeholder="Enter a topic to generate content..."
              onKeyDown={e => { if (e.key === 'Enter' && !e.shiftKey) { e.preventDefault(); handleGenerate(); } }}
              className="rounded-xl"
              data-testid="content-prompt-input"
            />
            <Button onClick={handleGenerate} disabled={loading || !prompt.trim()} className="rounded-xl px-6 btn-hover" data-testid="generate-content-btn">
              {loading ? <Loader2 className="h-4 w-4 animate-spin" /> : <Send className="h-4 w-4" />}
            </Button>
          </div>
        </div>
      </Card>

      {previewContent && (
        <ContentPreview content={previewContent} onClose={() => setPreviewContent(null)} onStatusUpdate={handleStatusUpdate} onDelete={handleDelete} onGenerateImage={handleGenerateImage} />
      )}
    </>
  );
}
