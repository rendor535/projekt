-- 1. Użytkownicy
CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     username VARCHAR(255) UNIQUE NOT NULL,
                                     password TEXT NOT NULL,
                                     role VARCHAR(20) NOT NULL CHECK (role IN ('trainer','trainee')),
                                     specialization VARCHAR(255),
                                     experience INTEGER NOT NULL DEFAULT 0,
                                     is_online BOOLEAN NOT NULL DEFAULT FALSE,
                                     created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                     last_seen TIMESTAMP DEFAULT NOW()
);

-- 2. Profile użytkowników
CREATE TABLE IF NOT EXISTS profiles (
                                        user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
                                        age INTEGER CHECK (age >= 12 AND age <= 100),
                                        height INTEGER CHECK (height >= 50 AND height <= 300),
                                        weight NUMERIC(5,2) CHECK (weight > 0),
                                        gender VARCHAR(10) CHECK (gender IN ('male','female','other')),
                                        goal VARCHAR(30) CHECK (
                                            goal IN (
                                                     'reduction',
                                                     'utrzymanie_wagi',
                                                     'przyrost_masy',
                                                     'budowanie_miesni',
                                                     'rekompozycja_ciala'
                                                )
                                            ),
                                        activity_level NUMERIC(3,2) NOT NULL DEFAULT 0,
                                        created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                        updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
--t--
DO $$
    BEGIN
        IF EXISTS (
            SELECT 1 FROM pg_trigger WHERE tgname = 'trigger_update_profiles_updated_at'
        ) THEN
            EXECUTE 'DROP TRIGGER trigger_update_profiles_updated_at ON profiles';
    END IF;
END$$;
-- 2a. Trigger do aktualizacji updated_at w profiles
CREATE OR REPLACE FUNCTION update_profiles_updated_at()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_profiles_updated_at
    BEFORE UPDATE ON profiles
    FOR EACH ROW EXECUTE FUNCTION update_profiles_updated_at();

-- 3. Plany treningowe
CREATE TABLE IF NOT EXISTS training_plans (
                                              id SERIAL PRIMARY KEY,
                                              client_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                              trainer_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                              name VARCHAR(255) NOT NULL,
                                              description TEXT,
                                              start_date DATE NOT NULL,
                                              end_date DATE,
                                              created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                              CHECK (end_date IS NULL OR start_date <= end_date)
);

-- 4. Treningi
CREATE TABLE IF NOT EXISTS exercise_library (
                                                id SERIAL PRIMARY KEY,
                                                name VARCHAR(255) UNIQUE NOT NULL,
                                                muscles TEXT NOT NULL,
                                                equipment TEXT,
                                                level VARCHAR(50) NOT NULL CHECK (level IN ('Początkujący', 'Średniozaawansowany', 'Zaawansowany'))
);
CREATE TABLE IF NOT EXISTS training_plan_days (
                                                  id SERIAL PRIMARY KEY,
                                                  training_plan_id INTEGER NOT NULL REFERENCES training_plans(id) ON DELETE CASCADE,
                                                  day_of_week VARCHAR(15) NOT NULL CHECK (
                                                      day_of_week IN (
                                                                      'poniedziałek', 'wtorek', 'środa', 'czwartek',
                                                                      'piątek', 'sobota', 'niedziela'
                                                          )
                                                      )
);
CREATE TABLE IF NOT EXISTS training_day_exercises (
                                                      id SERIAL PRIMARY KEY,
                                                      training_day_id INTEGER NOT NULL REFERENCES training_plan_days(id) ON DELETE CASCADE,
                                                      exercise_id INTEGER NOT NULL REFERENCES exercise_library(id) ON DELETE CASCADE,
                                                      sets INTEGER NOT NULL CHECK (sets >= 1),
                                                      reps INTEGER NOT NULL CHECK (reps >= 1),
                                                      weight INTEGER NOT NULL CHECK (weight >= 0),
                                                      instructions TEXT
);
-- 6. Pliki multimedialne
CREATE TABLE IF NOT EXISTS media_files (
                                           id SERIAL PRIMARY KEY,
                                           title VARCHAR(255) NOT NULL,
                                           url TEXT NOT NULL,
                                           uploaded_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
                                           created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 7. Przypisanie plików multimedialnych użytkownikom
CREATE TABLE IF NOT EXISTS user_media (
                                          id SERIAL PRIMARY KEY,
                                          uploader_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                          recipient_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                          file_name VARCHAR(255) NOT NULL,
                                          file_url VARCHAR(512) NOT NULL,
                                          file_type VARCHAR(50) NOT NULL,
                                          description TEXT,
                                          is_read BOOLEAN NOT NULL DEFAULT FALSE,
                                          created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 8. Plany dietetyczne
CREATE TABLE IF NOT EXISTS diet_plans (
                                          id SERIAL PRIMARY KEY,
                                          client_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                          trainer_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                          name VARCHAR(255) NOT NULL,
                                          description TEXT,
                                          calories INTEGER,
                                          meals TEXT,
                                          start_date DATE,
                                          end_date DATE,
                                          is_active BOOLEAN NOT NULL DEFAULT FALSE,
                                          created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 9. Posiłki w diecie
CREATE TABLE IF NOT EXISTS meals (
                                     id SERIAL PRIMARY KEY,
                                     diet_plan_id INTEGER NOT NULL REFERENCES diet_plans(id) ON DELETE CASCADE,
                                     meal_order INTEGER NOT NULL CHECK (meal_order BETWEEN 1 AND 5),
                                     name VARCHAR(100) NOT NULL,
                                     calories INTEGER NOT NULL CHECK (calories >= 0),
                                     proteins INTEGER NOT NULL CHECK (proteins >= 0),
                                     fats INTEGER NOT NULL CHECK (fats >= 0),
                                     carbs INTEGER NOT NULL CHECK (carbs >= 0)
);

-- 10. Produkty spożywcze
CREATE TABLE IF NOT EXISTS food_products (
                                             id SERIAL PRIMARY KEY,
                                             name VARCHAR(255) NOT NULL,
                                             calories NUMERIC(6,2) NOT NULL,
                                             protein NUMERIC(6,2) NOT NULL,
                                             fat NUMERIC(6,2) NOT NULL,
                                             carbs NUMERIC(6,2) NOT NULL
);

-- 11. Produkty w posiłkach
CREATE TABLE IF NOT EXISTS meal_products (
                                             id SERIAL PRIMARY KEY,
                                             meal_id INTEGER NOT NULL REFERENCES meals(id) ON DELETE CASCADE,
                                             product_id INTEGER NOT NULL REFERENCES food_products(id) ON DELETE CASCADE,
                                             amount_grams NUMERIC(6,2) NOT NULL CHECK (amount_grams > 0)
);

-- 12. Raporty postępów
CREATE TABLE IF NOT EXISTS progress_reports (
                                  id SERIAL PRIMARY KEY,
                                  client_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                  trainer_id INTEGER REFERENCES users(id) ON DELETE CASCADE,  -- NO NOT NULL!
                                  weight NUMERIC(5,2),
                                  body_fat NUMERIC(5,2),
                                  measurements JSONB,
                                  notes TEXT,
                                  photo_url VARCHAR(255),
                                  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- 13. Wiadomości (czat)
CREATE TABLE IF NOT EXISTS messages (
                                        id SERIAL PRIMARY KEY,
                                        sender_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                        receiver_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                        content TEXT NOT NULL,
                                        is_read BOOLEAN NOT NULL DEFAULT FALSE,
                                        created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 14. Trainer-Trainee
CREATE TABLE IF NOT EXISTS trainer_trainee (
                                               id SERIAL PRIMARY KEY,
                                               trainer_id INTEGER NOT NULL REFERENCES users(id),
                                               trainee_id INTEGER NOT NULL REFERENCES users(id),
                                               created_at TIMESTAMP NOT NULL,
                                               updated_at TIMESTAMP NOT NULL,
                                               UNIQUE(trainer_id, trainee_id)
);

-- 15. Logi wagi
CREATE TABLE IF NOT EXISTS weight_logs (
                                           id SERIAL PRIMARY KEY,
                                           user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                           weight NUMERIC(5,2) NOT NULL,
                                           body_fat NUMERIC(5,2),
                                           notes TEXT,
                                           created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 16. Ukończenia treningów
CREATE TABLE IF NOT EXISTS workout_completions (
                                                   id SERIAL PRIMARY KEY,
                                                   user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                                   exercise_assignment_id INTEGER NOT NULL REFERENCES training_day_exercises(id) ON DELETE CASCADE,
                                                   completed_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                                   duration_minutes INTEGER,
                                                   notes TEXT
);
CREATE TABLE IF NOT EXISTS exercise_videos (
                                               id SERIAL PRIMARY KEY,
                                               name VARCHAR(255) NOT NULL,
                                               youtube_url VARCHAR(512) NOT NULL,
                                               created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_workout_completion
    ON workout_completions(user_id, exercise_assignment_id, (DATE(completed_at)));

-- 17. Profile trenerów
CREATE TABLE IF NOT EXISTS trainer_profiles (
                                                user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
                                                first_name VARCHAR(100) NOT NULL,
                                                last_name VARCHAR(100) NOT NULL,
                                                bio TEXT NOT NULL,
                                                specializations JSONB NOT NULL DEFAULT '[]',
                                                experience INTEGER NOT NULL DEFAULT 0,
                                                certifications JSONB NOT NULL DEFAULT '[]',
                                                price_per_hour NUMERIC(10,2) NOT NULL DEFAULT 0,
                                                languages JSONB NOT NULL DEFAULT '[]',
                                                contact_email VARCHAR(255),
                                                contact_phone VARCHAR(20),
                                                location VARCHAR(255),
                                                avatar VARCHAR(500),
                                                is_active BOOLEAN NOT NULL DEFAULT TRUE,
                                                created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
                                                updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_trainer_profiles_active ON trainer_profiles(is_active);
CREATE INDEX IF NOT EXISTS idx_trainer_profiles_specializations ON trainer_profiles USING gin(specializations);
CREATE INDEX IF NOT EXISTS idx_trainer_profiles_experience ON trainer_profiles(experience);

-- 18. Żądania trenera
CREATE TABLE IF NOT EXISTS trainer_requests (
                                                id SERIAL PRIMARY KEY,
                                                client_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                                trainer_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                                status VARCHAR(20) NOT NULL DEFAULT 'pending',
                                                message TEXT,
                                                created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                                updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                                UNIQUE(client_id, trainer_id)
);
CREATE INDEX IF NOT EXISTS idx_trainer_requests_client ON trainer_requests(client_id);
CREATE INDEX IF NOT EXISTS idx_trainer_requests_trainer ON trainer_requests(trainer_id);

-- Indeksy pomocnicze
CREATE INDEX IF NOT EXISTS idx_messages_conversation ON messages(sender_id, receiver_id);
CREATE INDEX IF NOT EXISTS idx_messages_unread ON messages(receiver_id, is_read);
CREATE INDEX IF NOT EXISTS idx_progress_reports_client ON progress_reports(client_id);
CREATE INDEX IF NOT EXISTS idx_user_media_recipient ON user_media(recipient_id);
CREATE INDEX IF NOT EXISTS idx_diet_plans_client ON diet_plans(client_id);
-- Drop the existing unique constraint
ALTER TABLE trainer_requests DROP CONSTRAINT IF EXISTS trainer_requests_client_id_trainer_id_key;

-- Create a new unique constraint that only applies to pending requests
CREATE UNIQUE INDEX IF NOT EXISTS trainer_requests_client_trainer_pending_idx
    ON trainer_requests (client_id, trainer_id)
    WHERE status = 'pending';
-- Dodaj kolumny do śledzenia powiadomień
ALTER TABLE trainer_trainee
    ADD COLUMN IF NOT EXISTS last_reports_check TIMESTAMP DEFAULT NULL,
    ADD COLUMN IF NOT EXISTS last_media_check TIMESTAMP DEFAULT NULL,
    ADD COLUMN IF NOT EXISTS last_messages_check TIMESTAMP DEFAULT NULL;

-- Dodaj kolumny do diet_plans dla systemu opinii
ALTER TABLE diet_plans
    ADD COLUMN IF NOT EXISTS feedback_request BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS feedback_given BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS feedback_text TEXT;