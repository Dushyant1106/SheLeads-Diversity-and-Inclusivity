import { useState, useEffect } from 'react';
import api, { addUserIdToParams } from '@/lib/api';
import { toast } from 'sonner';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { DollarSign, TrendingDown, Percent, CheckCircle, RefreshCw } from 'lucide-react';

const CACHE_KEY = 'market_value_data';
const CACHE_EXPIRY_KEY = 'market_value_expiry';
const CACHE_DURATION = 24 * 60 * 60 * 1000; // 24 hours

export default function MarketValueLoans() {
  const [loading, setLoading] = useState(false);
  const [refreshing, setRefreshing] = useState(false);
  const [marketData, setMarketData] = useState(null);
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

        if (now < expiryTime) {
          const parsed = JSON.parse(cachedData);
          setMarketData(parsed.marketData);
          setLastUpdated(new Date(parsed.timestamp));
          return;
        }
      }
    } catch (error) {
      console.error('Failed to load cached data:', error);
    }
  };

  const fetchMarketData = async (isRefresh = false) => {
    try {
      if (isRefresh) {
        setRefreshing(true);
      } else {
        setLoading(true);
      }

      const response = await api.get('/analytics/market-value', {
        params: addUserIdToParams()
      });
      const marketDataResult = response.data.data;

      setMarketData(marketDataResult);

      // Cache the data
      const timestamp = Date.now();
      const cacheData = {
        marketData: marketDataResult,
        timestamp: timestamp
      };

      localStorage.setItem(CACHE_KEY, JSON.stringify(cacheData));
      localStorage.setItem(CACHE_EXPIRY_KEY, (timestamp + CACHE_DURATION).toString());
      setLastUpdated(new Date(timestamp));

      if (isRefresh) {
        toast.success('Market value refreshed!');
      }
    } catch (error) {
      toast.error('Failed to load market value data');
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  };

  const handleRefresh = () => {
    fetchMarketData(true);
  };

  if (loading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-48" />
        <Skeleton className="h-96" />
      </div>
    );
  }

  const totalMarketValue = marketData?.total_market_value || 0;
  const categoryBreakdown = marketData?.category_breakdown || [];
  const totalPoints = marketData?.total_points || 0;
  const loanOptions = marketData?.loan_options || [];
  const hasData = marketData !== null;

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
              <>No data yet. Click refresh to calculate your work value.</>
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
      {/* Market Value Overview */}
      <Card className="rounded-2xl bg-gradient-to-br from-green-50 to-emerald-50 dark:from-green-950/20 dark:to-emerald-950/20 border-green-200 dark:border-green-800">
        <CardHeader className="pb-3">
          <CardTitle className="text-lg flex items-center gap-2 heading-font">
            <DollarSign className="h-5 w-5 text-green-600" />
            Your Work's Market Value
          </CardTitle>
          <p className="text-xs text-muted-foreground mt-1">
            Based on professional service rates in your area
          </p>
        </CardHeader>
        <CardContent className="p-6">
          <div className="text-center mb-6">
            <p className="text-5xl font-bold heading-font text-green-600">
              ₹{totalMarketValue.toLocaleString('en-IN')}
            </p>
            <p className="text-sm text-muted-foreground mt-2">
              Total value of unpaid work logged (Indian market rates)
            </p>
          </div>

          {/* Category Breakdown */}
          <div className="space-y-3">
            <p className="text-sm font-medium text-muted-foreground uppercase tracking-wider">
              Breakdown by Category
            </p>
            {categoryBreakdown.map((item, index) => (
              <div
                key={index}
                className="flex items-center justify-between p-3 rounded-xl bg-white/50 dark:bg-gray-900/50 border"
              >
                <div className="flex items-center gap-3">
                  <div className="h-10 w-10 rounded-lg bg-green-100 dark:bg-green-900/30 flex items-center justify-center">
                    <DollarSign className="h-5 w-5 text-green-600" />
                  </div>
                  <div>
                    <p className="font-medium capitalize">{item.category}</p>
                    <p className="text-xs text-muted-foreground">
                      {item.hours.toFixed(1)} hours @ ₹{item.rate}/hr
                    </p>
                  </div>
                </div>
                <p className="font-bold text-green-600">
                  ₹{item.value.toLocaleString('en-IN')}
                </p>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Loan Options */}
      <Card className="rounded-2xl">
        <CardHeader className="pb-3 border-b">
          <CardTitle className="text-lg flex items-center gap-2 heading-font">
            <Percent className="h-5 w-5 text-primary" />
            Loan Options with Points Discount
          </CardTitle>
          <p className="text-xs text-muted-foreground mt-1">
            Your {totalPoints.toLocaleString()} points unlock reduced interest rates
          </p>
        </CardHeader>
        <CardContent className="p-6">
          {loanOptions.length === 0 ? (
            <div className="text-center py-8 text-muted-foreground">
              <TrendingDown className="h-12 w-12 mx-auto mb-3 opacity-50" />
              <p>Log more work to unlock loan benefits</p>
            </div>
          ) : (
            <div className="space-y-4">
              {loanOptions.map((loan, index) => (
                <div
                  key={index}
                  className="p-5 rounded-xl border bg-card hover:shadow-md transition-shadow"
                >
                  <div className="flex items-start justify-between mb-4">
                    <div>
                      <h3 className="font-semibold text-lg">{loan.name}</h3>
                      <p className="text-sm text-muted-foreground mt-1">
                        {loan.description}
                      </p>
                    </div>
                    <Badge variant="secondary" className="shrink-0">
                      {loan.term}
                    </Badge>
                  </div>

                  <div className="grid grid-cols-2 gap-4 mb-4">
                    <div className="p-3 rounded-lg bg-muted/50">
                      <p className="text-xs text-muted-foreground mb-1">Loan Amount</p>
                      <p className="text-xl font-bold heading-font">
                        ₹{loan.amount.toLocaleString('en-IN')}
                      </p>
                    </div>
                    <div className="p-3 rounded-lg bg-muted/50">
                      <p className="text-xs text-muted-foreground mb-1">Monthly EMI</p>
                      <p className="text-xl font-bold heading-font">
                        ₹{Math.round(loan.monthly_payment).toLocaleString('en-IN')}
                      </p>
                    </div>
                  </div>

                  {/* Interest Rate Comparison */}
                  <div className="flex items-center justify-between p-4 rounded-lg bg-gradient-to-r from-red-50 to-green-50 dark:from-red-950/20 dark:to-green-950/20">
                    <div className="flex items-center gap-4">
                      <div>
                        <p className="text-xs text-muted-foreground mb-1">Original Rate</p>
                        <p className="text-2xl font-bold text-red-600 line-through">
                          {loan.original_rate}%
                        </p>
                      </div>
                      <TrendingDown className="h-6 w-6 text-green-600" />
                      <div>
                        <p className="text-xs text-muted-foreground mb-1">Your Rate</p>
                        <p className="text-2xl font-bold text-green-600">
                          {loan.reduced_rate}%
                        </p>
                      </div>
                    </div>
                    <div className="text-right">
                      <Badge className="bg-green-600 text-white">
                        Save {loan.discount_percent.toFixed(1)}%
                      </Badge>
                      <p className="text-xs text-muted-foreground mt-1">
                        ₹{Math.round(loan.total_savings).toLocaleString('en-IN')} saved
                      </p>
                    </div>
                  </div>

                  {loan.eligible && (
                    <div className="flex items-center gap-2 mt-3 text-sm text-green-600">
                      <CheckCircle className="h-4 w-4" />
                      <span>You're eligible for this loan</span>
                    </div>
                  )}
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}

