import { useState, useEffect } from 'react';
import api, { addUserIdToParams } from '@/lib/api';
import { toast } from 'sonner';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { Calendar, Clock, Sparkles, Lightbulb, TrendingUp, RefreshCw } from 'lucide-react';

const CACHE_KEY = 'smart_calendar_data';
const CACHE_EXPIRY_KEY = 'smart_calendar_expiry';
const CACHE_DURATION = 24 * 60 * 60 * 1000; // 24 hours in milliseconds

export default function SmartCalendar() {
  const [loading, setLoading] = useState(false);
  const [refreshing, setRefreshing] = useState(false);
  const [calendarData, setCalendarData] = useState(null);
  const [businessProfile, setBusinessProfile] = useState(null);
  const [lastUpdated, setLastUpdated] = useState(null);

  useEffect(() => {
    loadCachedData();
  }, []);

  // Load cached data on mount
  const loadCachedData = () => {
    try {
      const cachedData = localStorage.getItem(CACHE_KEY);
      const cachedExpiry = localStorage.getItem(CACHE_EXPIRY_KEY);

      if (cachedData && cachedExpiry) {
        const expiryTime = parseInt(cachedExpiry, 10);
        const now = Date.now();

        // Check if cache is still valid
        if (now < expiryTime) {
          const parsed = JSON.parse(cachedData);
          setCalendarData(parsed.calendarData);
          setBusinessProfile(parsed.businessProfile);
          setLastUpdated(new Date(parsed.timestamp));
          return; // Data loaded from cache, no API call needed
        }
      }
    } catch (error) {
      console.error('Failed to load cached data:', error);
    }
  };

  // Fetch fresh data from API
  const fetchCalendarData = async (isRefresh = false) => {
    try {
      if (isRefresh) {
        setRefreshing(true);
      } else {
        setLoading(true);
      }

      // Fetch work pattern analysis
      const calendarResponse = await api.get('/analytics/calendar-insights', {
        params: addUserIdToParams()
      });
      const calendarDataResult = calendarResponse.data.data;

      // Fetch business profile for suggestions
      let businessProfileResult = null;
      try {
        const profileResponse = await api.get('/business/profile', {
          params: addUserIdToParams()
        });
        businessProfileResult = profileResponse.data.data;
      } catch (error) {
        // Business profile not set up yet
        businessProfileResult = null;
      }

      // Update state
      setCalendarData(calendarDataResult);
      setBusinessProfile(businessProfileResult);

      // Cache the data
      const timestamp = Date.now();
      const cacheData = {
        calendarData: calendarDataResult,
        businessProfile: businessProfileResult,
        timestamp: timestamp
      };

      localStorage.setItem(CACHE_KEY, JSON.stringify(cacheData));
      localStorage.setItem(CACHE_EXPIRY_KEY, (timestamp + CACHE_DURATION).toString());
      setLastUpdated(new Date(timestamp));

      if (isRefresh) {
        toast.success('Calendar insights refreshed!');
      }
    } catch (error) {
      toast.error('Failed to load calendar insights');
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  };

  // Handle refresh button click
  const handleRefresh = () => {
    fetchCalendarData(true);
  };

  if (loading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-48" />
        <Skeleton className="h-64" />
      </div>
    );
  }

  const freeTimeSlots = calendarData?.free_time_slots || [];
  const workPattern = calendarData?.work_pattern || {};
  const suggestions = calendarData?.business_suggestions || [];
  const hasData = calendarData !== null;

  // Format last updated time
  const formatLastUpdated = () => {
    if (!lastUpdated) return 'Never';
    const now = new Date();
    const diff = now - lastUpdated;
    const minutes = Math.floor(diff / 60000);
    const hours = Math.floor(diff / 3600000);
    const days = Math.floor(diff / 86400000);

    if (minutes < 1) return 'Just now';
    if (minutes < 60) return `${minutes}m ago`;
    if (hours < 24) return `${hours}h ago`;
    return `${days}d ago`;
  };

  return (
    <div className="space-y-6">
      {/* Header with Refresh Button */}
      <div className="flex items-center justify-between">
        <div>
          <p className="text-sm text-muted-foreground">
            {hasData ? (
              <>Last updated: {formatLastUpdated()}</>
            ) : (
              <>No data yet. Click refresh to analyze your work patterns.</>
            )}
          </p>
        </div>
        <Button
          onClick={handleRefresh}
          disabled={refreshing}
          variant="outline"
          size="sm"
          className="rounded-xl"
        >
          <RefreshCw className={`h-4 w-4 mr-2 ${refreshing ? 'animate-spin' : ''}`} />
          {refreshing ? 'Refreshing...' : 'Refresh'}
        </Button>
      </div>
      {/* Work Pattern Overview */}
      <Card className="rounded-2xl">
        <CardHeader className="pb-3 border-b">
          <CardTitle className="text-lg flex items-center gap-2 heading-font">
            <TrendingUp className="h-5 w-5 text-primary" />
            Your Work Pattern
          </CardTitle>
        </CardHeader>
        <CardContent className="p-6">
          <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
            <div className="text-center p-4 rounded-xl bg-muted/50">
              <p className="text-2xl font-bold heading-font text-primary">
                {workPattern.busiest_day || 'N/A'}
              </p>
              <p className="text-xs text-muted-foreground mt-1">Busiest Day</p>
            </div>
            <div className="text-center p-4 rounded-xl bg-muted/50">
              <p className="text-2xl font-bold heading-font text-secondary">
                {workPattern.avg_hours_per_day?.toFixed(1) || '0'}h
              </p>
              <p className="text-xs text-muted-foreground mt-1">Avg Hours/Day</p>
            </div>
            <div className="text-center p-4 rounded-xl bg-muted/50">
              <p className="text-2xl font-bold heading-font">
                {workPattern.peak_time || 'N/A'}
              </p>
              <p className="text-xs text-muted-foreground mt-1">Peak Time</p>
            </div>
            <div className="text-center p-4 rounded-xl bg-muted/50">
              <p className="text-2xl font-bold heading-font text-green-600">
                {workPattern.free_days || 0}
              </p>
              <p className="text-xs text-muted-foreground mt-1">Free hours/Week</p>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Free Time Slots */}
      <Card className="rounded-2xl">
        <CardHeader className="pb-3 border-b">
          <CardTitle className="text-lg flex items-center gap-2 heading-font">
            <Clock className="h-5 w-5 text-secondary" />
            Your Free Time Windows
          </CardTitle>
          <p className="text-xs text-muted-foreground mt-1">
            Based on your routine, these are your best times for business activities
          </p>
        </CardHeader>
        <CardContent className="p-6">
          {freeTimeSlots.length === 0 ? (
            <div className="text-center py-8 text-muted-foreground">
              <Calendar className="h-12 w-12 mx-auto mb-3 opacity-50" />
              <p>Log more work to see your free time patterns</p>
            </div>
          ) : (
            <div className="space-y-3">
              {freeTimeSlots.map((slot, index) => (
                <div
                  key={index}
                  className="flex items-center justify-between p-4 rounded-xl border bg-card hover:bg-muted/50 transition-colors"
                >
                  <div className="flex items-center gap-3">
                    <div className="h-10 w-10 rounded-lg bg-secondary/20 flex items-center justify-center">
                      <Clock className="h-5 w-5 text-secondary" />
                    </div>
                    <div>
                      <p className="font-medium">{slot.day}</p>
                      <p className="text-sm text-muted-foreground">{slot.time_range}</p>
                    </div>
                  </div>
                  <Badge variant="outline" className="capitalize">
                    {slot.duration}
                  </Badge>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Business Suggestions */}
      {businessProfile && suggestions.length > 0 && (
        <Card className="rounded-2xl border-primary/20 bg-primary/5">
          <CardHeader className="pb-3 border-b border-primary/20">
            <CardTitle className="text-lg flex items-center gap-2 heading-font">
              <Lightbulb className="h-5 w-5 text-primary" />
              Smart Business Suggestions
            </CardTitle>
            <p className="text-xs text-muted-foreground mt-1">
              Based on your free time and business goals
            </p>
          </CardHeader>
          <CardContent className="p-6">
            <div className="space-y-3">
              {suggestions.map((suggestion, index) => (
                <div
                  key={index}
                  className="flex items-start gap-3 p-4 rounded-xl bg-card border"
                >
                  <div className="h-8 w-8 rounded-lg bg-primary/20 flex items-center justify-center shrink-0 mt-0.5">
                    <Sparkles className="h-4 w-4 text-primary" />
                  </div>
                  <div className="flex-1">
                    <p className="font-medium text-sm">{suggestion.title}</p>
                    <p className="text-xs text-muted-foreground mt-1">
                      {suggestion.description}
                    </p>
                    <div className="flex items-center gap-2 mt-2">
                      <Badge variant="secondary" className="text-xs">
                        {suggestion.time_slot}
                      </Badge>
                      <Badge variant="outline" className="text-xs">
                        {suggestion.duration}
                      </Badge>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}

      {!businessProfile && (
        <Card className="rounded-2xl border-dashed">
          <CardContent className="p-8 text-center">
            <Lightbulb className="h-12 w-12 mx-auto mb-3 text-muted-foreground opacity-50" />
            <p className="font-medium mb-1">Set up your business profile</p>
            <p className="text-sm text-muted-foreground">
              Get personalized business activity suggestions based on your free time
            </p>
          </CardContent>
        </Card>
      )}
    </div>
  );
}


