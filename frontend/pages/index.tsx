import type { NextPage } from 'next'
import { useEffect, useState, useMemo, memo } from 'react'

// WailsのGoバインディングの型定義
declare global {
  interface Window {
    go: {
      main: {
        App: {
          GetRemainingTime: () => Promise<number> // ナノ秒を返す
          UpdateSessionState: () => Promise<void>
          GetCurrentSessionState: () => Promise<{
            mode: string
            current_session_started_at: string
            started_at: string
            paused: boolean
            session_count: number
          }>
          StartPomodoro: () => Promise<void>
          PausePomodoro: () => Promise<void>
          ResumePomodoro: () => Promise<void>
        }
      }
    }
  }
}

// 時間をフォーマットする関数（ナノ秒をミリ秒に変換して表示）
const formatTime = (nanoseconds: number): string => {
  const totalSeconds = Math.floor(nanoseconds / 1e9)
  const minutes = Math.floor(totalSeconds / 60)
  const seconds = totalSeconds % 60
  return `${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`
}

// タイマー表示コンポーネント（remainingTimeが変更されたときのみ再レンダリング）
const TimerDisplay = memo(({ remainingTime }: { remainingTime: number }) => {
  return (
    <div style={{
      fontSize: '120px',
      fontWeight: 'bold',
      color: 'white',
      textShadow: '2px 2px 4px rgba(0,0,0,0.3)',
      fontFamily: 'monospace',
    }}>
      {formatTime(remainingTime)}
    </div>
  )
})

TimerDisplay.displayName = 'TimerDisplay'

// モード表示コンポーネント（sessionStateが変更されたときのみ再レンダリング）
const ModeDisplay = memo(({ sessionState }: { sessionState: any }) => {
  if (!sessionState) return null

  const mode = sessionState.mode || sessionState.Mode
  const modeText = mode === 'work' ? '作業中' : 
                   mode === 'short_break' ? '短い休憩' : 
                   mode === 'long_break' ? '長い休憩' : ''

  return (
    <div style={{
      marginTop: '20px',
      fontSize: '24px',
      color: 'white',
      textShadow: '1px 1px 2px rgba(0,0,0,0.3)',
    }}>
      {modeText}
    </div>
  )
})

ModeDisplay.displayName = 'ModeDisplay'

// ボタンコンポーネント（sessionStateが変更されたときのみ再レンダリング）
const ControlButtons = memo(({ 
  isSessionStarted, 
  sessionState, 
  onStart, 
  onPause, 
  onResume 
}: { 
  isSessionStarted: boolean
  sessionState: any
  onStart: () => void
  onPause: () => void
  onResume: () => void
}) => {
  const buttonStyle = {
    padding: '15px 40px',
    fontSize: '20px',
    fontWeight: 'bold',
    color: 'white',
    backgroundColor: 'rgba(255, 255, 255, 0.2)',
    border: '2px solid white',
    borderRadius: '8px',
    cursor: 'pointer',
    transition: 'all 0.3s ease',
    backdropFilter: 'blur(10px)' as const,
  }

  const handleMouseOver = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.currentTarget.style.backgroundColor = 'rgba(255, 255, 255, 0.3)'
    e.currentTarget.style.transform = 'scale(1.05)'
  }

  const handleMouseOut = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.currentTarget.style.backgroundColor = 'rgba(255, 255, 255, 0.2)'
    e.currentTarget.style.transform = 'scale(1)'
  }

  if (!isSessionStarted) {
    return (
      <button
        onClick={onStart}
        style={buttonStyle}
        onMouseOver={handleMouseOver}
        onMouseOut={handleMouseOut}
      >
        開始
      </button>
    )
  }

  if (sessionState.paused || sessionState.Paused) {
    return (
      <button
        onClick={onResume}
        style={buttonStyle}
        onMouseOver={handleMouseOver}
        onMouseOut={handleMouseOut}
      >
        再開
      </button>
    )
  }

  return (
    <button
      onClick={onPause}
      style={buttonStyle}
      onMouseOver={handleMouseOver}
      onMouseOut={handleMouseOut}
    >
      一時停止
    </button>
  )
})

ControlButtons.displayName = 'ControlButtons'

const Home: NextPage = () => {
  const [remainingTime, setRemainingTime] = useState<number>(0)
  const [sessionState, setSessionState] = useState<any>(null)

  // セッション状態を更新する関数
  const updateSession = async () => {
    try {
      await window.go.main.App.UpdateSessionState()
      const state = await window.go.main.App.GetCurrentSessionState()
      setSessionState(state)
      const remaining = await window.go.main.App.GetRemainingTime()
      setRemainingTime(remaining)
      
      // タイマーが0になった場合は再読み込み
      if (remaining <= 0) {
        setTimeout(() => {
          window.location.reload()
        }, 100)
      }
    } catch (error) {
      console.error('[updateSession] Failed to update session:', error)
    }
  }

  // Startボタンのハンドラー
  const handleStart = async () => {
    try {
      await window.go.main.App.StartPomodoro()
      await updateSession()
    } catch (error) {
      console.error('[handleStart] Failed to start pomodoro:', error)
    }
  }

  // Pauseボタンのハンドラー
  const handlePause = async () => {
    try {
      await window.go.main.App.PausePomodoro()
      await updateSession()
    } catch (error) {
      console.error('[handlePause] Failed to pause pomodoro:', error)
    }
  }

  // Resumeボタンのハンドラー
  const handleResume = async () => {
    try {
      await window.go.main.App.ResumePomodoro()
      await updateSession()
    } catch (error) {
      console.error('[handleResume] Failed to resume pomodoro:', error)
    }
  }

  // セッションが開始されているかどうかを判定
  const isSessionStarted = useMemo(() => {
    const currentModeStartedAt = sessionState?.CurrentModeStartedAt || sessionState?.current_session_started_at || ''
    return sessionState && 
      currentModeStartedAt && 
      currentModeStartedAt !== '' &&
      currentModeStartedAt !== '0001-01-01T00:00:00Z' &&
      currentModeStartedAt !== '0001-01-01T00:00:00.000Z'
  }, [sessionState])

  useEffect(() => {
    // 初回読み込み時にセッション状態を更新
    updateSession()

    // 1秒ごとにタイマーを更新
    const interval = setInterval(async () => {
      try {
        const remaining = await window.go.main.App.GetRemainingTime()
        setRemainingTime(remaining)
        
        // タイマーが0になった場合は再読み込み
        if (remaining <= 0) {
          clearInterval(interval)
          setTimeout(() => {
            window.location.reload()
          }, 100)
        }
      } catch (error) {
        console.error('[useEffect] Failed to get remaining time:', error)
      }
    }, 1000)

    // ウィンドウのフォーカス時にセッション状態を更新
    const handleFocus = () => {
      updateSession()
    }

    window.addEventListener('focus', handleFocus)

    // クリーンアップ
    return () => {
      clearInterval(interval)
      window.removeEventListener('focus', handleFocus)
    }
  }, [])

  return (
    <div style={{
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'center',
      height: '100vh',
      width: '100vw',
    }}>
      <div style={{
        textAlign: 'center',
      }}>
        <TimerDisplay remainingTime={remainingTime} />
        <ModeDisplay sessionState={sessionState} />
        <div style={{
          marginTop: '40px',
          display: 'flex',
          gap: '20px',
          justifyContent: 'center',
        }}>
          <ControlButtons
            isSessionStarted={isSessionStarted}
            sessionState={sessionState}
            onStart={handleStart}
            onPause={handlePause}
            onResume={handleResume}
          />
        </div>
      </div>
    </div>
  )
}

export default Home
