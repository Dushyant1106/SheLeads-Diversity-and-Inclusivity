import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '@/lib/auth-context';
import api, { addUserIdToParams } from '@/lib/api';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { ClipboardList, Megaphone, Trophy, Clock, FileText, TrendingUp, ArrowRight, Sparkles } from 'lucide-react';
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, Cell } from 'recharts';

const CHART_COLORS = ['#E07A5F', '#81B29A', '#F2CC8F', '#3D405B', '#E63946'];

export default function DashboardPage() {
  const { user } = useAuth();
  const navigate = useNavigate();
  const [analytics, setAnalytics] = useState(null);
  const [metrics, setMetrics] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [analyticsRes, metricsRes] = await Promise.allSettled([
          api.get('/analytics/summary', { params: addUserIdToParams() }),
          api.get('/metrics/aggregated', { params: addUserIdToParams() }),
        ]);
        if (analyticsRes.status === 'fulfilled') setAnalytics(analyticsRes.value.data?.data);
        if (metricsRes.status === 'fulfilled') setMetrics(metricsRes.value.data?.data);
      } catch {
        /* silently fail */
      }
      setLoading(false);
    };
    fetchData();
  }, []);

  const statCards = [
    { label: 'Points Earned', value: analytics?.total_points ?? 0, icon: Trophy, color: 'text-primary' },
    { label: 'Hours Logged', value: analytics?.total_hours ? analytics.total_hours.toFixed(1) : '0', icon: Clock, color: 'text-secondary' },
    { label: 'Content Created', value: metrics?.total_posts ?? 0, icon: FileText, color: 'text-accent-foreground' },
    { label: 'Engagement Rate', value: metrics?.avg_engagement_rate ? `${(metrics.avg_engagement_rate * 100).toFixed(1)}%` : '0%', icon: TrendingUp, color: 'text-primary' },
  ];

  const categoryData = analytics?.by_category?.map((cat, i) => ({
    name: cat.category || cat._id || `Cat ${i + 1}`,
    hours: cat.total_hours || cat.hours || 0,
  })) || [];

  return (
    <div className="space-y-8" data-testid="dashboard-page">
      {/* Welcome Section */}
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-3xl sm:text-4xl font-bold tracking-tight heading-font">
            Welcome, {user?.name?.split(' ')[0] || 'there'}
          </h1>
          <p className="text-muted-foreground mt-1 text-sm sm:text-base">Here's your overview for today</p>
        </div>
        <div className="flex gap-2">
          <Button onClick={() => navigate('/logger')} className="rounded-xl btn-hover" data-testid="quick-log-work-btn">
            <ClipboardList className="h-4 w-4 mr-2" /> Log Work
          </Button>
          <Button variant="outline" onClick={() => navigate('/marketing')} className="rounded-xl btn-hover" data-testid="quick-create-content-btn">
            <Megaphone className="h-4 w-4 mr-2" /> Create Content
          </Button>
        </div>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 animate-stagger" data-testid="stats-grid">
        {statCards.map(({ label, value, icon: Icon, color }) => (
          <Card key={label} className="rounded-2xl border hover:shadow-[0_8px_30px_rgb(0,0,0,0.12)] transition-shadow">
            <CardContent className="p-5">
              {loading ? (
                <Skeleton className="h-16 w-full" />
              ) : (
                <>
                  <div className="flex items-center justify-between mb-3">
                    <span className="text-xs uppercase tracking-widest text-muted-foreground font-medium">{label}</span>
                    <Icon className={`h-4 w-4 ${color}`} />
                  </div>
                  <p className="text-2xl sm:text-3xl font-bold tracking-tight heading-font">{value}</p>
                </>
              )}
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Charts Row */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Category Hours Chart */}
        <Card className="rounded-2xl">
          <CardHeader className="pb-2">
            <CardTitle className="text-lg heading-font">Hours by Category</CardTitle>
          </CardHeader>
          <CardContent>
            {loading ? (
              <Skeleton className="h-48 w-full" />
            ) : categoryData.length > 0 ? (
              <ResponsiveContainer width="100%" height={200}>
                <BarChart data={categoryData}>
                  <XAxis dataKey="name" tick={{ fontSize: 11 }} />
                  <YAxis tick={{ fontSize: 11 }} />
                  <Tooltip />
                  <Bar dataKey="hours" radius={[6, 6, 0, 0]}>
                    {categoryData.map((entry, i) => (
                      <Cell key={`cell-${i}`} fill={CHART_COLORS[i % CHART_COLORS.length]} />
                    ))}
                  </Bar>
                </BarChart>
              </ResponsiveContainer>
            ) : (
              <div className="h-48 flex flex-col items-center justify-center text-muted-foreground">
                <Sparkles className="h-8 w-8 mb-2 opacity-50" />
                <p className="text-sm">No data yet. Start logging your work!</p>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Quick Actions */}
        <Card className="rounded-2xl">
          <CardHeader className="pb-2">
            <CardTitle className="text-lg heading-font">Quick Actions</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            {[
              { label: 'Log Care Work', desc: 'Track your daily contributions', to: '/logger', icon: ClipboardList },
              { label: 'Generate Content', desc: 'Create blogs and social posts', to: '/marketing', icon: Megaphone },
              { label: 'Download Report', desc: 'Get your monthly PDF report', to: '/logger', icon: FileText },
            ].map(({ label, desc, to, icon: Icon }) => (
              <button
                key={label}
                onClick={() => navigate(to)}
                className="w-full flex items-center gap-4 p-4 rounded-xl bg-muted/30 hover:bg-muted/60 transition-colors text-left group"
                data-testid={`quick-action-${label.toLowerCase().replace(/\s+/g, '-')}`}
              >
                <div className="h-10 w-10 rounded-lg bg-primary/10 flex items-center justify-center shrink-0">
                  <Icon className="h-5 w-5 text-primary" />
                </div>
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-medium">{label}</p>
                  <p className="text-xs text-muted-foreground">{desc}</p>
                </div>
                <ArrowRight className="h-4 w-4 text-muted-foreground group-hover:text-foreground transition-colors" />
              </button>
            ))}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
