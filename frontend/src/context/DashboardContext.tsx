import React, { createContext, useContext, useReducer, useCallback, useEffect } from 'react'
import {
  Dashboard,
  DashboardWidget,
  DashboardState,
  DashboardAction,
  DashboardContextType,
  CreateDashboardRequest,
  UpdateDashboardRequest,
  CreateWidgetRequest,
  UpdateWidgetRequest,
  WidgetDataResponse,
  ExportOptions,
  ExportResult,
  BIError
} from '../types/business-intelligence'
import { biService } from '../services/business-intelligence'

// Initial state
const initialState: DashboardState = {
  currentDashboard: null,
  dashboards: [],
  widgets: [],
  isLoading: false,
  isEditing: false,
  selectedWidget: null,
  draggedWidget: null,
  showWidgetPalette: false,
  showSettings: false,
  error: null
}

// Reducer function
function dashboardReducer(state: DashboardState, action: DashboardAction): DashboardState {
  switch (action.type) {
    case 'SET_CURRENT_DASHBOARD':
      return {
        ...state,
        currentDashboard: action.payload,
        widgets: action.payload.widgets || []
      }
    
    case 'SET_DASHBOARDS':
      return {
        ...state,
        dashboards: action.payload
      }
    
    case 'ADD_DASHBOARD':
      return {
        ...state,
        dashboards: [...state.dashboards, action.payload]
      }
    
    case 'UPDATE_DASHBOARD':
      return {
        ...state,
        dashboards: state.dashboards.map(d => 
          d.id === action.payload.id ? action.payload : d
        ),
        currentDashboard: state.currentDashboard?.id === action.payload.id 
          ? action.payload 
          : state.currentDashboard
      }
    
    case 'REMOVE_DASHBOARD':
      return {
        ...state,
        dashboards: state.dashboards.filter(d => d.id !== action.payload),
        currentDashboard: state.currentDashboard?.id === action.payload 
          ? null 
          : state.currentDashboard
      }
    
    case 'ADD_WIDGET':
      const newWidgets = [...state.widgets, action.payload]
      return {
        ...state,
        widgets: newWidgets,
        currentDashboard: state.currentDashboard 
          ? { ...state.currentDashboard, widgets: newWidgets }
          : null
      }
    
    case 'UPDATE_WIDGET':
      const updatedWidgets = state.widgets.map(w => 
        w.id === action.payload.id ? action.payload : w
      )
      return {
        ...state,
        widgets: updatedWidgets,
        currentDashboard: state.currentDashboard 
          ? { ...state.currentDashboard, widgets: updatedWidgets }
          : null,
        selectedWidget: state.selectedWidget?.id === action.payload.id 
          ? action.payload 
          : state.selectedWidget
      }
    
    case 'REMOVE_WIDGET':
      const filteredWidgets = state.widgets.filter(w => w.id !== action.payload)
      return {
        ...state,
        widgets: filteredWidgets,
        currentDashboard: state.currentDashboard 
          ? { ...state.currentDashboard, widgets: filteredWidgets }
          : null,
        selectedWidget: state.selectedWidget?.id === action.payload 
          ? null 
          : state.selectedWidget
      }
    
    case 'SET_EDITING':
      return {
        ...state,
        isEditing: action.payload,
        selectedWidget: action.payload ? state.selectedWidget : null,
        showWidgetPalette: action.payload ? state.showWidgetPalette : false
      }
    
    case 'SET_SELECTED_WIDGET':
      return {
        ...state,
        selectedWidget: action.payload
      }
    
    case 'SET_DRAGGED_WIDGET':
      return {
        ...state,
        draggedWidget: action.payload
      }
    
    case 'SET_WIDGET_PALETTE_VISIBLE':
      return {
        ...state,
        showWidgetPalette: action.payload
      }
    
    case 'SET_SETTINGS_VISIBLE':
      return {
        ...state,
        showSettings: action.payload
      }
    
    case 'SET_LOADING':
      return {
        ...state,
        isLoading: action.payload
      }
    
    case 'SET_ERROR':
      return {
        ...state,
        error: action.payload,
        isLoading: false
      }
    
    default:
      return state
  }
}

// Create context
const DashboardContext = createContext<DashboardContextType | undefined>(undefined)

// Provider component
export function DashboardProvider({ children }: { children: React.ReactNode }) {
  const [state, dispatch] = useReducer(dashboardReducer, initialState)

  // Load dashboards on mount
  useEffect(() => {
    loadDashboards()
  }, [])

  // Dashboard operations
  const loadDashboards = useCallback(async () => {
    try {
      dispatch({ type: 'SET_LOADING', payload: true })
      const response = await biService.getDashboards()
      dispatch({ type: 'SET_DASHBOARDS', payload: response.dashboards.map(d => d.dashboard) })
      
      // Set first dashboard as current if none selected
      if (!state.currentDashboard && response.dashboards.length > 0) {
        const firstDashboard = response.dashboards[0].dashboard
        dispatch({ type: 'SET_CURRENT_DASHBOARD', payload: firstDashboard })
      }
    } catch (error) {
      dispatch({ 
        type: 'SET_ERROR', 
        payload: { 
          type: 'dashboard', 
          message: 'Failed to load dashboards',
          details: error
        }
      })
    } finally {
      dispatch({ type: 'SET_LOADING', payload: false })
    }
  }, [state.currentDashboard])

  const createDashboard = useCallback(async (request: CreateDashboardRequest): Promise<Dashboard> => {
    try {
      dispatch({ type: 'SET_LOADING', payload: true })
      const dashboard = await biService.createDashboard(request)
      dispatch({ type: 'ADD_DASHBOARD', payload: dashboard })
      
      // Set as current dashboard if it's the default
      if (request.isDefault || state.dashboards.length === 0) {
        dispatch({ type: 'SET_CURRENT_DASHBOARD', payload: dashboard })
      }
      
      return dashboard
    } catch (error) {
      dispatch({ 
        type: 'SET_ERROR', 
        payload: { 
          type: 'dashboard', 
          message: 'Failed to create dashboard',
          details: error
        }
      })
      throw error
    } finally {
      dispatch({ type: 'SET_LOADING', payload: false })
    }
  }, [state.dashboards.length])

  const updateDashboard = useCallback(async (id: string, request: UpdateDashboardRequest): Promise<Dashboard> => {
    try {
      const dashboard = await biService.updateDashboard(id, request)
      dispatch({ type: 'UPDATE_DASHBOARD', payload: dashboard })
      return dashboard
    } catch (error) {
      dispatch({ 
        type: 'SET_ERROR', 
        payload: { 
          type: 'dashboard', 
          message: 'Failed to update dashboard',
          details: error
        }
      })
      throw error
    }
  }, [])

  const deleteDashboard = useCallback(async (id: string): Promise<void> => {
    try {
      await biService.deleteDashboard(id)
      dispatch({ type: 'REMOVE_DASHBOARD', payload: id })
      
      // If we deleted the current dashboard, select another one
      if (state.currentDashboard?.id === id && state.dashboards.length > 1) {
        const remainingDashboards = state.dashboards.filter(d => d.id !== id)
        dispatch({ type: 'SET_CURRENT_DASHBOARD', payload: remainingDashboards[0] })
      }
    } catch (error) {
      dispatch({ 
        type: 'SET_ERROR', 
        payload: { 
          type: 'dashboard', 
          message: 'Failed to delete dashboard',
          details: error
        }
      })
      throw error
    }
  }, [state.currentDashboard, state.dashboards])

  const loadDashboard = useCallback(async (id: string): Promise<void> => {
    try {
      dispatch({ type: 'SET_LOADING', payload: true })
      const dashboard = await biService.getDashboard(id)
      dispatch({ type: 'SET_CURRENT_DASHBOARD', payload: dashboard })
    } catch (error) {
      dispatch({ 
        type: 'SET_ERROR', 
        payload: { 
          type: 'dashboard', 
          message: 'Failed to load dashboard',
          details: error
        }
      })
    } finally {
      dispatch({ type: 'SET_LOADING', payload: false })
    }
  }, [])

  // Widget operations
  const createWidget = useCallback(async (request: CreateWidgetRequest): Promise<DashboardWidget> => {
    try {
      const widget = await biService.createWidget(request)
      dispatch({ type: 'ADD_WIDGET', payload: widget })
      return widget
    } catch (error) {
      dispatch({ 
        type: 'SET_ERROR', 
        payload: { 
          type: 'widget', 
          message: 'Failed to create widget',
          details: error
        }
      })
      throw error
    }
  }, [])

  const updateWidget = useCallback(async (id: string, request: UpdateWidgetRequest): Promise<DashboardWidget> => {
    try {
      const widget = await biService.updateWidget(id, request)
      dispatch({ type: 'UPDATE_WIDGET', payload: widget })
      return widget
    } catch (error) {
      dispatch({ 
        type: 'SET_ERROR', 
        payload: { 
          type: 'widget', 
          message: 'Failed to update widget',
          details: error
        }
      })
      throw error
    }
  }, [])

  const deleteWidget = useCallback(async (id: string): Promise<void> => {
    try {
      await biService.deleteWidget(id)
      dispatch({ type: 'REMOVE_WIDGET', payload: id })
    } catch (error) {
      dispatch({ 
        type: 'SET_ERROR', 
        payload: { 
          type: 'widget', 
          message: 'Failed to delete widget',
          details: error
        }
      })
      throw error
    }
  }, [])

  const getWidgetData = useCallback(async (id: string): Promise<WidgetDataResponse> => {
    try {
      return await biService.getWidgetData(id)
    } catch (error) {
      dispatch({ 
        type: 'SET_ERROR', 
        payload: { 
          type: 'data', 
          message: 'Failed to load widget data',
          details: error
        }
      })
      throw error
    }
  }, [])

  // Export operations
  const exportDashboard = useCallback(async (id: string, options: ExportOptions): Promise<ExportResult> => {
    try {
      return await biService.exportDashboard(id, options)
    } catch (error) {
      dispatch({ 
        type: 'SET_ERROR', 
        payload: { 
          type: 'dashboard', 
          message: 'Failed to export dashboard',
          details: error
        }
      })
      throw error
    }
  }, [])

  // Widget manipulation helpers
  const moveWidget = useCallback((widgetId: string, newPosition: { x: number; y: number }) => {
    const widget = state.widgets.find(w => w.id === widgetId)
    if (widget) {
      const updatedWidget = {
        ...widget,
        position: { ...widget.position, ...newPosition }
      }
      dispatch({ type: 'UPDATE_WIDGET', payload: updatedWidget })
      
      // Async update to backend
      updateWidget(widgetId, { position: updatedWidget.position }).catch(error => {
        console.error('Failed to update widget position:', error)
        // Revert the optimistic update
        dispatch({ type: 'UPDATE_WIDGET', payload: widget })
      })
    }
  }, [state.widgets, updateWidget])

  const resizeWidget = useCallback((widgetId: string, newSize: { width: number; height: number }) => {
    const widget = state.widgets.find(w => w.id === widgetId)
    if (widget) {
      const updatedWidget = {
        ...widget,
        size: { ...widget.size, ...newSize }
      }
      dispatch({ type: 'UPDATE_WIDGET', payload: updatedWidget })
      
      // Async update to backend
      updateWidget(widgetId, { size: updatedWidget.size }).catch(error => {
        console.error('Failed to update widget size:', error)
        // Revert the optimistic update
        dispatch({ type: 'UPDATE_WIDGET', payload: widget })
      })
    }
  }, [state.widgets, updateWidget])

  // UI state helpers
  const toggleEditing = useCallback(() => {
    dispatch({ type: 'SET_EDITING', payload: !state.isEditing })
  }, [state.isEditing])

  const selectWidget = useCallback((widget: DashboardWidget | null) => {
    dispatch({ type: 'SET_SELECTED_WIDGET', payload: widget })
  }, [])

  const toggleWidgetPalette = useCallback(() => {
    dispatch({ type: 'SET_WIDGET_PALETTE_VISIBLE', payload: !state.showWidgetPalette })
  }, [state.showWidgetPalette])

  const toggleSettings = useCallback(() => {
    dispatch({ type: 'SET_SETTINGS_VISIBLE', payload: !state.showSettings })
  }, [state.showSettings])

  const clearError = useCallback(() => {
    dispatch({ type: 'SET_ERROR', payload: null })
  }, [])

  // Bulk operations
  const bulkUpdateWidgets = useCallback(async (updates: Array<{ id: string; request: UpdateWidgetRequest }>) => {
    try {
      const results = await biService.bulkUpdateWidgets(updates)
      results.forEach(widget => {
        dispatch({ type: 'UPDATE_WIDGET', payload: widget })
      })
      return results
    } catch (error) {
      dispatch({ 
        type: 'SET_ERROR', 
        payload: { 
          type: 'widget', 
          message: 'Failed to update widgets',
          details: error
        }
      })
      throw error
    }
  }, [])

  const bulkDeleteWidgets = useCallback(async (ids: string[]) => {
    try {
      await biService.bulkDeleteWidgets(ids)
      ids.forEach(id => {
        dispatch({ type: 'REMOVE_WIDGET', payload: id })
      })
    } catch (error) {
      dispatch({ 
        type: 'SET_ERROR', 
        payload: { 
          type: 'widget', 
          message: 'Failed to delete widgets',
          details: error
        }
      })
      throw error
    }
  }, [])

  // Context value
  const contextValue: DashboardContextType = {
    state,
    dispatch,
    
    // Dashboard operations
    createDashboard,
    updateDashboard,
    deleteDashboard,
    loadDashboard,
    loadDashboards,
    
    // Widget operations
    createWidget,
    updateWidget,
    deleteWidget,
    getWidgetData,
    moveWidget,
    resizeWidget,
    
    // Export operations
    exportDashboard,
    
    // UI state
    toggleEditing,
    selectWidget,
    toggleWidgetPalette,
    toggleSettings,
    clearError,
    
    // Bulk operations
    bulkUpdateWidgets,
    bulkDeleteWidgets
  }

  return (
    <DashboardContext.Provider value={contextValue}>
      {children}
    </DashboardContext.Provider>
  )
}

// Hook to use dashboard context
export function useDashboard(): DashboardContextType {
  const context = useContext(DashboardContext)
  if (context === undefined) {
    throw new Error('useDashboard must be used within a DashboardProvider')
  }
  return context
}

export default DashboardContext