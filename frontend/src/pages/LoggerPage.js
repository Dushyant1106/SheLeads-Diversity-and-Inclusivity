import { useState, useEffect } from 'react';
import api, { addUserIdToParams, addUserIdToFormData } from '@/lib/api';
import { toast } from 'sonner';
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Badge } from '@/components/ui/badge';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Skeleton } from '@/components/ui/skeleton';
import {
  Upload, Loader2, Trophy, Clock, AlertTriangle, CheckCircle,
  Download, BarChart3, Sparkles
} from 'lucide-react';
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip } from 'recharts';
import SmartCalendar from '@/components/logger/SmartCalendar';
import MarketValueLoans from '@/components/logger/MarketValueLoans';

const CATEGORIES = [
  { id: 'cooking', label: 'Cooking', img: 'https://images.unsplash.com/photo-1758523421916-8721f7532953?w=200&h=200&fit=crop' },
  { id: 'cleaning', label: 'Cleaning', img: 'https://images.unsplash.com/photo-1758273238370-3bc08e399620?w=200&h=200&fit=crop' },
  { id: 'childcare', label: 'Childcare', img: 'https://images.unsplash.com/photo-1758598738260-4893695cead9?w=200&h=200&fit=crop' },
  { id: 'eldercare', label: 'Eldercare', img: 'https://images.pexels.com/photos/7551667/pexels-photo-7551667.jpeg?auto=compress&w=200&h=200&fit=crop' },
  { id: 'laundry', label: 'Laundry', img: 'https://images.pexels.com/photos/5591640/pexels-photo-5591640.jpeg?auto=compress&w=200&h=200&fit=crop' },
  { id: 'shopping', label: 'Shopping', img: 'https://images.unsplash.com/photo-1583561552285-b38038c88e06?w=200&h=200&fit=crop' },
  { id: 'gardening', label: 'Gardening', img: 'https://images.unsplash.com/photo-1720517380910-1b9b209a4284?w=200&h=200&fit=crop' },
  { id: 'home_maintenance', label: 'Home Maint.', img: 'https://images.unsplash.com/photo-1581578731548-c64695cc6952?w=200&h=200&fit=crop' },
  { id: 'other', label: 'Other', img: '' },
];

const CHART_COLORS = ['#E07A5F', '#81B29A', '#F2CC8F', '#3D405B', '#E63946', '#8B5CF6', '#06B6D4', '#F59E0B', '#94A3B8'];

function WorkLogTab() {
  const [category, setCategory] = useState('');
  const [hours, setHours] = useState('');
  const [description, setDescription] = useState('');
  const [image, setImage] = useState(null);
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState(null);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!category || !hours || !image) {
      toast.error('Please fill all required fields and upload an image');
      return;
    }
    setLoading(true);
    setResult(null);
    try {
      const formData = new FormData();
      formData.append('category', category);
      formData.append('hours', parseFloat(hours));
      formData.append('description', description);
      formData.append('image', image);
      addUserIdToFormData(formData);
      const response = await api.post('/work/log', formData);
      setResult(response.data.data);
      toast.success(`Logged ${hours} hours of ${category}!`);
      setCategory('');
      setHours('');
      setDescription('');
      setImage(null);
    } catch (error) {
      toast.error(error.response?.data?.message || 'Failed to log work');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      {result && (
        <Card className="rounded-2xl border-secondary/50 bg-secondary/5 animate-fade-in" data-testid="log-result-card">
          <CardContent className="p-5 flex flex-col sm:flex-row items-start sm:items-center gap-4">
            <div className="h-12 w-12 rounded-full bg-secondary/20 flex items-center justify-center shrink-0">
              <Trophy className="h-6 w-6 text-secondary" />
            </div>
            <div className="flex-1">
              <p className="font-semibold text-lg heading-font">
                +{result.points_earned || 0} Points Earned!
              </p>
              {result.ai_verification && (
                <div className="flex items-center gap-2 mt-1">
                  <CheckCircle className="h-4 w-4 text-secondary" />
                  <span className="text-sm text-muted-foreground">AI Verified</span>
                </div>
              )}
              {result.burnout_alert && (
                <div className="flex items-center gap-2 mt-1">
                  <AlertTriangle className="h-4 w-4 text-destructive" />
                  <span className="text-sm text-destructive">Burnout alert: Consider taking a break</span>
                </div>
              )}
            </div>
          </CardContent>
        </Card>
      )}

      <div>
        <Label className="text-sm font-medium mb-3 block uppercase tracking-widest text-muted-foreground">Select Category</Label>
        <div className="grid grid-cols-3 sm:grid-cols-4 lg:grid-cols-5 gap-3 animate-stagger" data-testid="category-grid">
          {CATEGORIES.map(({ id, label, img }) => (
            <button
              key={id}
              type="button"
              onClick={() => setCategory(id)}
              data-testid={`category-${id}`}
              className={`category-card relative overflow-hidden rounded-2xl border-2 p-2 sm:p-3 text-center transition-all ${
                category === id
                  ? 'border-primary bg-primary/5 shadow-md'
                  : 'border-transparent bg-muted/30 hover:bg-muted/50'
              }`}
            >
              {img ? (
                <div className="w-full aspect-square rounded-xl overflow-hidden mb-2">
                  <img src={img} alt={label} className="w-full h-full object-cover" loading="lazy" />
                </div>
              ) : (
                <div className="w-full aspect-square rounded-xl bg-muted/50 flex items-center justify-center mb-2">
                  <Sparkles className="h-6 w-6 text-muted-foreground" />
                </div>
              )}
              <span className="text-xs font-medium">{label}</span>
            </button>
          ))}
        </div>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div className="space-y-2">
          <Label htmlFor="hours">Hours Spent *</Label>
          <Input id="hours" type="number" min="0.1" max="24" step="0.1" value={hours} onChange={e => setHours(e.target.value)} placeholder="e.g., 2.5" required data-testid="hours-input" className="rounded-xl" />
        </div>
        <div className="space-y-2">
          <Label htmlFor="work-image">Upload Image *</Label>
          <Input id="work-image" type="file" accept="image/*" onChange={e => setImage(e.target.files?.[0] || null)} className="cursor-pointer rounded-xl" data-testid="image-upload-input" />
          {image && <p className="text-xs text-muted-foreground truncate">{image.name}</p>}
        </div>
      </div>

      <div className="space-y-2">
        <Label htmlFor="description">Description (Optional)</Label>
        <Textarea id="description" value={description} onChange={e => setDescription(e.target.value)} placeholder="Describe your work..." rows={3} data-testid="description-input" className="rounded-xl" />
      </div>

      <Button onClick={handleSubmit} disabled={loading || !category || !hours || !image} className="rounded-xl h-11 px-8 btn-hover" data-testid="log-work-submit-btn">
        {loading ? <Loader2 className="h-4 w-4 animate-spin mr-2" /> : <Upload className="h-4 w-4 mr-2" />}
        Log Work Entry
      </Button>
    </div>
  );
}

function AnalyticsTab() {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.get('/analytics/summary', { params: addUserIdToParams() })
      .then(res => setData(res.data?.data))
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <div className="space-y-4"><Skeleton className="h-32" /><Skeleton className="h-64" /></div>;

  const categoryData = data?.by_category?.map((cat, i) => ({
    name: cat.category || cat._id || `Cat ${i + 1}`,
    hours: cat.total_hours || cat.hours || 0,
  })) || [];

  return (
    <div className="space-y-6" data-testid="analytics-tab">
      <div className="grid grid-cols-2 gap-4">
        <Card className="rounded-2xl">
          <CardContent className="p-5 text-center">
            <Clock className="h-6 w-6 text-secondary mx-auto mb-2" />
            <p className="text-3xl font-bold heading-font">{data?.total_hours?.toFixed(1) || 0}</p>
            <p className="text-xs uppercase tracking-widest text-muted-foreground mt-1">Total Hours</p>
          </CardContent>
        </Card>
        <Card className="rounded-2xl">
          <CardContent className="p-5 text-center">
            <Trophy className="h-6 w-6 text-primary mx-auto mb-2" />
            <p className="text-3xl font-bold heading-font">{data?.total_points || 0}</p>
            <p className="text-xs uppercase tracking-widest text-muted-foreground mt-1">Total Points</p>
          </CardContent>
        </Card>
      </div>

      {categoryData.length > 0 ? (
        <Card className="rounded-2xl">
          <CardHeader><CardTitle className="heading-font">Hours by Category</CardTitle></CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie data={categoryData} dataKey="hours" nameKey="name" cx="50%" cy="50%" outerRadius={100} label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}>
                  {categoryData.map((entry, i) => <Cell key={`cell-${i}`} fill={CHART_COLORS[i % CHART_COLORS.length]} />)}
                </Pie>
                <Tooltip />
              </PieChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>
      ) : (
        <Card className="rounded-2xl">
          <CardContent className="p-12 text-center text-muted-foreground">
            <BarChart3 className="h-10 w-10 mx-auto mb-3 opacity-50" />
            <p>No analytics data yet. Start logging your work to see insights!</p>
          </CardContent>
        </Card>
      )}
    </div>
  );
}

function ReportsTab() {
  const [month, setMonth] = useState(String(new Date().getMonth() + 1));
  const [year, setYear] = useState(String(new Date().getFullYear()));
  const [loading, setLoading] = useState(false);

  const handleDownload = async () => {
    setLoading(true);
    try {
      const response = await api.get('/reports/monthly/pdf', {
        params: addUserIdToParams({ month: parseInt(month), year: parseInt(year) }),
        responseType: 'blob',
      });
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', `sheleads_report_${year}_${month}.pdf`);
      document.body.appendChild(link);
      link.click();
      link.remove();
      window.URL.revokeObjectURL(url);
      toast.success('Report downloaded!');
    } catch (error) {
      toast.error(error.response?.data?.message || 'Failed to download report');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Card className="rounded-2xl max-w-lg" data-testid="reports-tab">
      <CardHeader>
        <CardTitle className="heading-font">Monthly Report</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <p className="text-sm text-muted-foreground">Download your monthly PDF report summarizing all logged work, points earned, and analytics.</p>
        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-2">
            <Label>Month</Label>
            <Select value={month} onValueChange={setMonth}>
              <SelectTrigger data-testid="report-month-select" className="rounded-xl"><SelectValue /></SelectTrigger>
              <SelectContent>
                {Array.from({ length: 12 }, (_, i) => (
                  <SelectItem key={i + 1} value={String(i + 1)}>
                    {new Date(2024, i).toLocaleString('default', { month: 'long' })}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <div className="space-y-2">
            <Label>Year</Label>
            <Select value={year} onValueChange={setYear}>
              <SelectTrigger data-testid="report-year-select" className="rounded-xl"><SelectValue /></SelectTrigger>
              <SelectContent>
                {[2024, 2025, 2026].map(y => (
                  <SelectItem key={y} value={String(y)}>{y}</SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </div>
        <Button onClick={handleDownload} disabled={loading} className="rounded-xl btn-hover" data-testid="download-report-btn">
          {loading ? <Loader2 className="h-4 w-4 animate-spin mr-2" /> : <Download className="h-4 w-4 mr-2" />}
          Download PDF Report
        </Button>
      </CardContent>
    </Card>
  );
}

export default function LoggerPage() {
  return (
    <div className="space-y-6" data-testid="logger-page">
      <div>
        <h1 className="text-3xl sm:text-4xl font-bold tracking-tight heading-font">
          Work Logger
        </h1>
        <p className="text-muted-foreground mt-1">Track and value your contributions</p>
      </div>

      <Tabs defaultValue="log" className="w-full">
        <TabsList className="rounded-xl">
          <TabsTrigger value="log" className="rounded-lg" data-testid="tab-log-work">Log Work</TabsTrigger>
          <TabsTrigger value="analytics" className="rounded-lg" data-testid="tab-analytics">Analytics</TabsTrigger>
          <TabsTrigger value="calendar" className="rounded-lg" data-testid="tab-calendar">Smart Calendar</TabsTrigger>
          <TabsTrigger value="market" className="rounded-lg" data-testid="tab-market">Market Value</TabsTrigger>
          <TabsTrigger value="reports" className="rounded-lg" data-testid="tab-reports">Reports</TabsTrigger>
        </TabsList>
        <TabsContent value="log" className="mt-6"><WorkLogTab /></TabsContent>
        <TabsContent value="analytics" className="mt-6"><AnalyticsTab /></TabsContent>
        <TabsContent value="calendar" className="mt-6"><SmartCalendar /></TabsContent>
        <TabsContent value="market" className="mt-6"><MarketValueLoans /></TabsContent>
        <TabsContent value="reports" className="mt-6"><ReportsTab /></TabsContent>
      </Tabs>
    </div>
  );
}
