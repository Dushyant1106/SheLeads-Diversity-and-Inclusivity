import { useState } from 'react';
import api, { addUserIdToParams } from '@/lib/api';
import { toast } from 'sonner';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Badge } from '@/components/ui/badge';
import { Loader2, Plus, X, Building2 } from 'lucide-react';

export default function BusinessProfileForm({ profile, onSuccess, isEdit = false }) {
  const [loading, setLoading] = useState(false);
  const [form, setForm] = useState({
    business_name: profile?.business_name || '',
    industry: profile?.industry || '',
    location: profile?.location || '',
    description: profile?.description || '',
    target_audience: profile?.target_audience || '',
    website: profile?.website || '',
    unique_selling_points: profile?.unique_selling_points || [],
    social_media_handles: profile?.social_media_handles || {},
  });
  const [uspInput, setUspInput] = useState('');
  const [socialPlatform, setSocialPlatform] = useState('');
  const [socialHandle, setSocialHandle] = useState('');

  const handleChange = (e) => {
    setForm(prev => ({ ...prev, [e.target.name]: e.target.value }));
  };

  const addUSP = () => {
    if (uspInput.trim()) {
      setForm(prev => ({ ...prev, unique_selling_points: [...prev.unique_selling_points, uspInput.trim()] }));
      setUspInput('');
    }
  };

  const removeUSP = (index) => {
    setForm(prev => ({ ...prev, unique_selling_points: prev.unique_selling_points.filter((_, i) => i !== index) }));
  };

  const addSocialHandle = () => {
    if (socialPlatform.trim() && socialHandle.trim()) {
      setForm(prev => ({
        ...prev,
        social_media_handles: { ...prev.social_media_handles, [socialPlatform.toLowerCase()]: socialHandle },
      }));
      setSocialPlatform('');
      setSocialHandle('');
    }
  };

  const removeSocialHandle = (platform) => {
    setForm(prev => {
      const handles = { ...prev.social_media_handles };
      delete handles[platform];
      return { ...prev, social_media_handles: handles };
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!form.business_name || !form.industry) {
      toast.error('Business name and industry are required');
      return;
    }
    setLoading(true);
    try {
      const response = await api.post('/business/profile', form, { params: addUserIdToParams() });
      toast.success(isEdit ? 'Profile updated!' : 'Business profile created!');
      onSuccess?.(response.data.data);
    } catch (error) {
      toast.error(error.response?.data?.message || 'Failed to save profile');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Card className="rounded-2xl max-w-2xl" data-testid="business-profile-form">
      <CardHeader>
        <div className="flex items-center gap-3">
          <div className="h-10 w-10 rounded-xl bg-primary/10 flex items-center justify-center">
            <Building2 className="h-5 w-5 text-primary" />
          </div>
          <div>
            <CardTitle className="heading-font">{isEdit ? 'Edit Business Profile' : 'Setup Business Profile'}</CardTitle>
            <CardDescription>Tell us about your business to personalize content</CardDescription>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-5">
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="business_name">Business Name *</Label>
              <Input id="business_name" name="business_name" value={form.business_name} onChange={handleChange} placeholder="EmpowerHer Crafts" required data-testid="business-name-input" className="rounded-xl" />
            </div>
            <div className="space-y-2">
              <Label htmlFor="industry">Industry *</Label>
              <Input id="industry" name="industry" value={form.industry} onChange={handleChange} placeholder="Handmade Crafts & Jewelry" required data-testid="industry-input" className="rounded-xl" />
            </div>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="location">Location</Label>
              <Input id="location" name="location" value={form.location} onChange={handleChange} placeholder="Mumbai, India" data-testid="location-input" className="rounded-xl" />
            </div>
            <div className="space-y-2">
              <Label htmlFor="website">Website</Label>
              <Input id="website" name="website" value={form.website} onChange={handleChange} placeholder="https://yourbusiness.com" data-testid="website-input" className="rounded-xl" />
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="bp-description">Business Description</Label>
            <Textarea id="bp-description" name="description" value={form.description} onChange={handleChange} placeholder="Describe your business..." rows={3} data-testid="bp-description-input" className="rounded-xl" />
          </div>

          <div className="space-y-2">
            <Label htmlFor="target_audience">Target Audience</Label>
            <Input id="target_audience" name="target_audience" value={form.target_audience} onChange={handleChange} placeholder="Women aged 25-45 who appreciate handmade products" data-testid="target-audience-input" className="rounded-xl" />
          </div>

          <div className="space-y-2">
            <Label>Unique Selling Points</Label>
            <div className="flex gap-2">
              <Input value={uspInput} onChange={e => setUspInput(e.target.value)} placeholder="Add a unique selling point" onKeyDown={e => { if (e.key === 'Enter') { e.preventDefault(); addUSP(); } }} data-testid="usp-input" className="rounded-xl" />
              <Button type="button" variant="outline" size="icon" onClick={addUSP} data-testid="add-usp-btn" className="shrink-0 rounded-xl"><Plus className="h-4 w-4" /></Button>
            </div>
            {form.unique_selling_points.length > 0 && (
              <div className="flex flex-wrap gap-2 mt-2">
                {form.unique_selling_points.map((usp, i) => (
                  <Badge key={i} variant="secondary" className="pl-3 pr-1 py-1 gap-1 rounded-lg">
                    {usp}
                    <button type="button" onClick={() => removeUSP(i)} className="ml-1 hover:text-destructive"><X className="h-3 w-3" /></button>
                  </Badge>
                ))}
              </div>
            )}
          </div>

          <div className="space-y-2">
            <Label>Social Media Handles</Label>
            <div className="flex gap-2">
              <Input value={socialPlatform} onChange={e => setSocialPlatform(e.target.value)} placeholder="Platform" className="w-1/3 rounded-xl" data-testid="social-platform-input" />
              <Input value={socialHandle} onChange={e => setSocialHandle(e.target.value)} placeholder="@handle" className="flex-1 rounded-xl" data-testid="social-handle-input" />
              <Button type="button" variant="outline" size="icon" onClick={addSocialHandle} data-testid="add-social-btn" className="shrink-0 rounded-xl"><Plus className="h-4 w-4" /></Button>
            </div>
            {Object.entries(form.social_media_handles).length > 0 && (
              <div className="flex flex-wrap gap-2 mt-2">
                {Object.entries(form.social_media_handles).map(([platform, handle]) => (
                  <Badge key={platform} variant="outline" className="pl-3 pr-1 py-1 gap-1 rounded-lg">
                    {platform}: {handle}
                    <button type="button" onClick={() => removeSocialHandle(platform)} className="ml-1 hover:text-destructive"><X className="h-3 w-3" /></button>
                  </Badge>
                ))}
              </div>
            )}
          </div>

          <Button type="submit" disabled={loading} className="rounded-xl h-11 px-8 btn-hover" data-testid="save-profile-btn">
            {loading && <Loader2 className="h-4 w-4 animate-spin mr-2" />}
            {isEdit ? 'Update Profile' : 'Create Profile'}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}
