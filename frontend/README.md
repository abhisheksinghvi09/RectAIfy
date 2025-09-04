# RectAIfy Web Interface

A modern, elegant React + MUI web interface for the RectAIfy startup idea validation platform.

## 🎨 Design Philosophy

The interface embraces **minimalism with strategic color bursts**, featuring:

- **Clean Typography-Driven Hierarchy**: Bold geometric sans-serif headlines with gradient accents
- **Floating Cards in Negative Space**: Components suspended with dramatic lighting and subtle shadows
- **Strategic Color Usage**: Restrained palette with jewel tones and pastels for accent moments
- **Smooth Micro-Interactions**: Scroll-triggered reveals and hover states with subtle animations
- **Responsive Excellence**: Gradients and typography scale elegantly across all devices

## 🚀 Features

### Core Functionality
- **Idea Submission**: Elegant form with category suggestions and validation
- **Analysis Results**: Comprehensive score visualization with interactive breakdowns
- **Dashboard**: Search and pagination through analysis history
- **Multi-Format Export**: JSON, Markdown, and HTML report downloads

### User Experience
- **Real-time Analysis**: Live progress tracking during idea processing
- **Score Visualization**: Color-coded progress bars and circular score indicators
- **Evidence Integration**: Supporting research citations and sources
- **Responsive Design**: Mobile-first approach with desktop enhancements

## 🛠️ Tech Stack

- **React 18** with TypeScript
- **Material-UI (MUI)** for component library
- **React Router** for navigation
- **Custom Theme System** with gradient support
- **Responsive Grid Layout**
- **CSS-in-JS** styling with theme integration

## 📦 Installation

1. **Navigate to the web directory:**
   ```bash
   cd web
   ```

2. **Install dependencies:**
   ```bash
   npm install
   ```

3. **Configure environment:**
   ```bash
   cp .env.example .env
   ```
   
   Edit `.env` to set your API URL:
   ```
   VITE_API_URL=http://localhost:9444
   ```

4. **Start development server:**
   ```bash
   npm run dev
   ```

The interface will be available at `http://localhost:5173`

## 🎯 Usage

### Submitting Ideas
1. Navigate to the home page
2. Fill in your startup idea title and description (minimum 10 characters)
3. Optionally select a category and location
4. Click "Analyze Idea" to submit

### Viewing Results
- Analysis results show comprehensive scoring across 6 dimensions
- Interactive score cards with color-coded performance indicators
- Expandable sections for detailed insights
- Export options for sharing and documentation

### Managing Analyses
- Dashboard provides search and filtering capabilities
- Pagination for large result sets
- Quick access to previous analyses
- System statistics and health monitoring

## 🎨 Theme Customization

The theme system supports extensive customization:

```typescript
// Custom gradient definitions
gradient: {
  primary: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
  secondary: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
  tertiary: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
}

// Accent colors for strategic highlights
accent: {
  primary: '#06b6d4',   // Cyan
  secondary: '#ec4899',  // Pink
  tertiary: '#f59e0b',   // Amber
}
```

## 📱 Responsive Breakpoints

- **Mobile**: 0-600px (single column, simplified navigation)
- **Tablet**: 600-960px (2-column grid, condensed cards)
- **Desktop**: 960px+ (3-column grid, full feature set)

## 🔧 Development

### Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint

### Project Structure

```
src/
├── components/          # Reusable UI components
│   ├── Layout.tsx      # Main app layout with navigation
│   └── ScoreCard.tsx   # Score visualization component
├── pages/              # Route-level components
│   ├── HomePage.tsx    # Idea submission form
│   ├── AnalysisPage.tsx # Results display
│   ├── DashboardPage.tsx # Analysis history
│   └── AboutPage.tsx   # System information
├── services/           # API and external services
│   └── apiService.ts   # RectAIfy API client
├── theme/              # MUI theme configuration
│   └── theme.ts        # Custom theme with gradients
├── types/              # TypeScript definitions
│   └── api.ts          # API response types
└── App.tsx             # Root component with routing
```

## 🌐 API Integration

The interface connects to the RectAIfy API with the following endpoints:

- `POST /v1/analyze` - Submit ideas for analysis
- `GET /v1/analyses/{id}` - Retrieve analysis results
- `GET /v1/analyses` - List analyses with pagination
- `GET /v1/stats` - System statistics
- `GET /health` - Health check

## 🎯 Performance Optimizations

- **Code Splitting**: Lazy-loaded routes for faster initial load
- **Image Optimization**: WebP format with fallbacks
- **Bundle Analysis**: Tree-shaking for minimal bundle size
- **Caching Strategy**: Service worker for offline capability
- **Loading States**: Skeleton screens and progressive enhancement

## 🔐 Security

- **Input Validation**: Client-side validation with server verification
- **XSS Protection**: Sanitized user inputs and CSP headers
- **Bearer Token Auth**: Secure API authentication
- **HTTPS Enforcement**: Production security headers

## 📈 Analytics & Monitoring

Integration points for analytics:

- **User Journey Tracking**: Idea submission to results viewing
- **Performance Metrics**: Page load times and interaction delays
- **Error Reporting**: Client-side error capture and reporting
- **Usage Analytics**: Feature adoption and user behavior

## 🤝 Contributing

1. Follow the existing code style and conventions
2. Use TypeScript for all new components
3. Include responsive design for all features
4. Test across mobile and desktop breakpoints
5. Follow the design system color and spacing guidelines

## 📄 License

This project is part of the RectAIfy platform. See the main project license for details.
