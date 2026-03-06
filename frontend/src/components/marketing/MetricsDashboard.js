import { useState, useEffect } from 'react';
import api, { addUserIdToParams } from '@/lib/api';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { Heart, MessageCircle, Share2, Eye, Users, TrendingUp, BarChart3 } from 'lucide-react';
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, Cell } from 'recharts';

const CHART_COLORS = ['#E07A5F', '#81B29A', '#F2CC8F', '#3D405B', '#E63946'];

export default function MetricsDashboard() {
  const [metrics, setMetrics] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.get('/metrics/aggregated', { params: addUserIdToParams() })
      .then(res => setMetrics(res.data?.data))
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  const statCards = [
    { label: 'Total Posts', value: metrics?.total_posts ?? 0, icon: BarChart3, color: 'text-primary' },
    { label: 'Likes', value: metrics?.total_likes ?? 0, icon: Heart, color: 'text-destructive' },
    { label: 'Comments', value: metrics?.total_comments ?? 0, icon: MessageCircle, color: 'text-secondary' },
    { label: 'Shares', value: metrics?.total_shares ?? 0, icon: Share2, color: 'text-primary' },
    { label: 'Impressions', value: metrics?.total_impressions ?? 0, icon: Eye, color: 'text-accent-foreground' },
    { label: 'Reach', value: metrics?.total_reach ?? 0, icon: Users, color: 'text-secondary' },
  ];

  const engagementRate = metrics?.avg_engagement_rate ? (metrics.avg_engagement_rate * 100).toFixed(2) : '0';

  const platformData = metrics?.by_platform?.map((p, i) => ({
    name: p.platform || p._id || `Platform ${i + 1}`,
    posts: p.total_posts || p.count || 0,
  })) || [];

  return (
    <div className="space-y-6" data-testid="metrics-dashboard">
      {/* Engagement Rate Hero */}
      <Card className="rounded-2xl bg-primary/5 border-primary/20">
        <CardContent className="p-6 flex items-center gap-4">
          <div className="h-16 w-16 rounded-2xl bg-primary/10 flex items-center justify-center">
            <TrendingUp className="h-8 w-8 text-primary" />
          </div>
          <div>
            <p className="text-xs uppercase tracking-widest text-muted-foreground">Avg Engagement Rate</p>
            {loading ? <Skeleton className="h-10 w-24 mt-1" /> : (
              <p className="text-4xl font-bold text-primary heading-font">{engagementRate}%</p>
            )}
          </div>
        </CardContent>
      </Card>

      {/* Stats Grid */}
      <div className="grid grid-cols-2 md:grid-cols-3 gap-4 animate-stagger">
        {statCards.map(({ label, value, icon: Icon, color }) => (
          <Card key={label} className="rounded-2xl">
            <CardContent className="p-4">
              {loading ? <Skeleton className="h-16" /> : (
                <>
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-[10px] uppercase tracking-widest text-muted-foreground">{label}</span>
                    <Icon className={`h-4 w-4 ${color}`} />
                  </div>
                  <p className="text-2xl font-bold heading-font">{typeof value === 'number' ? value.toLocaleString() : value}</p>
                </>
              )}
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Platform Breakdown */}
      {platformData.length > 0 && (
        <Card className="rounded-2xl">
          <CardHeader><CardTitle className="heading-font">Posts by Platform</CardTitle></CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={250}>
              <BarChart data={platformData}>
                <XAxis dataKey="name" tick={{ fontSize: 11 }} />
                <YAxis tick={{ fontSize: 11 }} />
                <Tooltip />
                <Bar dataKey="posts" radius={[6, 6, 0, 0]}>
                  {platformData.map((entry, i) => <Cell key={`cell-${i}`} fill={CHART_COLORS[i % CHART_COLORS.length]} />)}
                </Bar>
              </BarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>
      )}

      {!loading && !metrics && (
        <Card className="rounded-2xl">
          <CardContent className="p-12 text-center text-muted-foreground">
            <BarChart3 className="h-10 w-10 mx-auto mb-3 opacity-50" />
            <p>No metrics data yet. Post content and track performance!</p>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
