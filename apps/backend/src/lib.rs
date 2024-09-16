pub mod core {
    use axum::{
        routing::get,
        Router,
    };

    pub fn define_app_routes() -> axum::Router {
        let app_router = Router::new()
            .route("/health", get(health_handler));

        return app_router;
    }

    //  ____   ___  _   _ _____ _____ ____  
    // |  _ \ / _ \| | | |_   _| ____/ ___| 
    // | |_) | | | | | | | | | |  _| \___ \ 
    // |  _ <| |_| | |_| | | | | |___ ___) |
    // |_| \_\\___/ \___/  |_| |_____|____/ 
    //                                      

    async fn health_handler() -> &'static str {
        return "woof ! taking your thoughts to the pound";
    }
}
