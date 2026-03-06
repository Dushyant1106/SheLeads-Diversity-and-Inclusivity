import { useState, useEffect, useCallback } from 'react';
import api, { addUserIdToParams } from '@/lib/api';
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs';
import { Skeleton } from '@/components/ui/skeleton';
import BusinessProfileForm from '@/components/marketing/BusinessProfileForm';
import ContentChat from '@/components/marketing/ContentChat';
import ContentLibrary from '@/components/marketing/ContentLibrary';
import MetricsDashboard from '@/components/marketing/MetricsDashboard';
import BrandAssets from '@/components/marketing/BrandAssets';

export default function MarketingPage() {
  const [businessProfile, setBusinessProfile] = useState(null);
  const [hasProfile, setHasProfile] = useState(false);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState('content');

  const fetchProfile = useCallback(async () => {
    try {
      const response = await api.get('/business/profile', { params: addUserIdToParams() });
      if (response.data?.data) {
        setBusinessProfile(response.data.data);
        setHasProfile(true);
      }
    } catch (error) {
      if (error.response?.status === 404) {
        setHasProfile(false);
      }
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => { fetchProfile(); }, [fetchProfile]);

  const handleProfileSaved = (profile) => {
    setBusinessProfile(profile);
    setHasProfile(true);
  };

  if (loading) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-10 w-64" />
        <Skeleton className="h-96 w-full" />
      </div>
    );
  }

  if (!hasProfile) {
    return (
      <div className="space-y-6" data-testid="marketing-page-setup">
        <div>
          <h1 className="text-3xl sm:text-4xl font-bold tracking-tight heading-font">
            Setup Your Business
          </h1>
          <p className="text-muted-foreground mt-1">Let's get your business profile ready for marketing</p>
        </div>
        <BusinessProfileForm onSuccess={handleProfileSaved} />
      </div>
    );
  }

  return (
    <div className="space-y-6" data-testid="marketing-page">
      <div>
        <h1 className="text-3xl sm:text-4xl font-bold tracking-tight heading-font">
          Marketing Hub
        </h1>
        <p className="text-muted-foreground mt-1">{businessProfile?.business_name || 'Your Business'}</p>
      </div>

      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="rounded-xl">
          <TabsTrigger value="content" className="rounded-lg" data-testid="tab-content-studio">Content Studio</TabsTrigger>
          <TabsTrigger value="assets" className="rounded-lg" data-testid="tab-assets">Brand Assets</TabsTrigger>
          <TabsTrigger value="metrics" className="rounded-lg" data-testid="tab-metrics">Metrics</TabsTrigger>
          <TabsTrigger value="settings" className="rounded-lg" data-testid="tab-settings">Settings</TabsTrigger>
        </TabsList>

        <TabsContent value="content" className="mt-6">
          <div className="grid grid-cols-1 lg:grid-cols-5 gap-6">
            <div className="lg:col-span-3">
              <ContentChat />
            </div>
            <div className="lg:col-span-2">
              <ContentLibrary />
            </div>
          </div>
        </TabsContent>

        <TabsContent value="assets" className="mt-6">
          <div className="max-w-2xl">
            <BrandAssets />
          </div>
        </TabsContent>

        <TabsContent value="metrics" className="mt-6">
          <MetricsDashboard />
        </TabsContent>

        <TabsContent value="settings" className="mt-6">
          <BusinessProfileForm profile={businessProfile} onSuccess={handleProfileSaved} isEdit />
        </TabsContent>
      </Tabs>
    </div>
  );
}
