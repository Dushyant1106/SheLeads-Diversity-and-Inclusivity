import { useState, useEffect } from 'react';
import api, { addUserIdToParams } from '@/lib/api';
import { toast } from 'sonner';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Upload, Trash2, Image as ImageIcon, Loader2 } from 'lucide-react';

export default function BrandAssets() {
  const [assets, setAssets] = useState([]);
  const [logoUrl, setLogoUrl] = useState(null);
  const [loading, setLoading] = useState(true);
  const [uploading, setUploading] = useState(false);

  useEffect(() => {
    fetchAssets();
  }, []);

  const fetchAssets = async () => {
    try {
      setLoading(true);
      const response = await api.get('/assets', { params: addUserIdToParams() });
      const data = response.data.data;
      
      setLogoUrl(data.logo_url);
      
      // Flatten all asset types into a single array
      const allAssets = [];
      if (data.assets) {
        Object.entries(data.assets).forEach(([type, assetList]) => {
          if (Array.isArray(assetList)) {
            assetList.forEach(asset => {
              allAssets.push({ ...asset, type });
            });
          }
        });
      }
      
      // Sort by upload date (newest first)
      allAssets.sort((a, b) => new Date(b.uploaded_at) - new Date(a.uploaded_at));
      setAssets(allAssets);
    } catch (error) {
      toast.error('Failed to load assets');
    } finally {
      setLoading(false);
    }
  };

  const handleFileUpload = async (event, assetType = 'reference_image') => {
    const file = event.target.files?.[0];
    if (!file) return;

    // Validate file type
    const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp'];
    if (!allowedTypes.includes(file.type)) {
      toast.error('Please upload an image file (JPG, PNG, GIF, or WebP)');
      return;
    }

    // Validate file size (10MB)
    if (file.size > 10 * 1024 * 1024) {
      toast.error('File size must be less than 10MB');
      return;
    }

    try {
      setUploading(true);
      const formData = new FormData();
      formData.append('file', file);

      const params = addUserIdToParams();
      params.asset_type = assetType;

      await api.post('/assets/upload', formData, {
        params,
        headers: { 'Content-Type': 'multipart/form-data' },
      });

      toast.success(`${assetType === 'logo' ? 'Logo' : 'Asset'} uploaded successfully!`);
      await fetchAssets();
    } catch (error) {
      toast.error('Failed to upload: ' + (error.response?.data?.error || 'Unknown error'));
    } finally {
      setUploading(false);
      event.target.value = ''; // Reset input
    }
  };

  if (loading) {
    return (
      <Card className="rounded-2xl">
        <CardHeader className="pb-3 border-b">
          <CardTitle className="text-lg flex items-center gap-2 heading-font">
            <ImageIcon className="h-5 w-5 text-primary" />
            Brand Assets
          </CardTitle>
        </CardHeader>
        <CardContent className="p-6">
          <div className="flex items-center justify-center h-40">
            <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="rounded-2xl">
      <CardHeader className="pb-3 border-b">
        <CardTitle className="text-lg flex items-center gap-2 heading-font">
          <ImageIcon className="h-5 w-5 text-primary" />
          Brand Assets
        </CardTitle>
        <p className="text-xs text-muted-foreground mt-1">
          Upload your logo, product images, and brand visuals. These will be used as references for AI image generation.
        </p>
      </CardHeader>
      <CardContent className="p-6 space-y-6">
        {/* Logo Upload */}
        <div className="space-y-2">
          <label className="text-sm font-medium">Logo</label>
          <div className="flex items-center gap-3">
            {logoUrl ? (
              <div className="relative w-20 h-20 rounded-lg border overflow-hidden bg-muted">
                <img src={logoUrl} alt="Logo" className="w-full h-full object-contain" />
              </div>
            ) : (
              <div className="w-20 h-20 rounded-lg border border-dashed flex items-center justify-center bg-muted/50">
                <ImageIcon className="h-6 w-6 text-muted-foreground" />
              </div>
            )}
            <div>
              <input
                type="file"
                id="logo-upload"
                className="hidden"
                accept="image/*"
                onChange={(e) => handleFileUpload(e, 'logo')}
                disabled={uploading}
              />
              <Button
                size="sm"
                variant="outline"
                className="rounded-lg"
                onClick={() => document.getElementById('logo-upload').click()}
                disabled={uploading}
              >
                <Upload className="h-3 w-3 mr-1" />
                {logoUrl ? 'Change Logo' : 'Upload Logo'}
              </Button>
            </div>
          </div>
        </div>

        {/* Reference Images */}
        <div className="space-y-3">
          <div className="flex items-center justify-between">
            <label className="text-sm font-medium">Reference Images ({assets.length})</label>
            <div>
              <input
                type="file"
                id="asset-upload"
                className="hidden"
                accept="image/*"
                onChange={(e) => handleFileUpload(e, 'reference_image')}
                disabled={uploading}
              />
              <Button
                size="sm"
                className="rounded-lg h-8"
                onClick={() => document.getElementById('asset-upload').click()}
                disabled={uploading}
              >
                {uploading ? (
                  <Loader2 className="h-3 w-3 mr-1 animate-spin" />
                ) : (
                  <Upload className="h-3 w-3 mr-1" />
                )}
                Upload Image
              </Button>
            </div>
          </div>

          {/* Asset Grid */}
          {assets.length === 0 ? (
            <div className="border border-dashed rounded-lg p-8 text-center">
              <ImageIcon className="h-8 w-8 mx-auto text-muted-foreground mb-2" />
              <p className="text-sm text-muted-foreground">No assets uploaded yet</p>
              <p className="text-xs text-muted-foreground mt-1">
                Upload product images, brand photos, or reference visuals
              </p>
            </div>
          ) : (
            <div className="grid grid-cols-3 gap-3 max-h-96 overflow-y-auto custom-scrollbar">
              {assets.map((asset, index) => (
                <div
                  key={index}
                  className="relative group aspect-square rounded-lg border overflow-hidden bg-muted hover:ring-2 hover:ring-primary transition-all"
                >
                  <img
                    src={asset.url}
                    alt={asset.filename}
                    className="w-full h-full object-cover"
                  />
                  <div className="absolute inset-0 bg-black/60 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
                    <div className="text-center text-white p-2">
                      <p className="text-xs font-medium truncate">{asset.filename}</p>
                      <p className="text-[10px] text-white/70 mt-1">
                        {new Date(asset.uploaded_at).toLocaleDateString()}
                      </p>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        <div className="text-xs text-muted-foreground bg-muted/50 rounded-lg p-3">
          <p className="font-medium mb-1">💡 Tip:</p>
          <p>Upload 3-5 high-quality images of your products, workspace, or brand style. The AI will use these as references to generate images that match your brand identity.</p>
        </div>
      </CardContent>
    </Card>
  );
}


