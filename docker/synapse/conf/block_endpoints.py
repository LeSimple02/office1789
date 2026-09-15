"""
Module pour bloquer les endpoints dangereux dans Synapse
"""
from synapse.http.servlet import RestServlet
from synapse.api.errors import Codes, SynapseError

class BlockPasswordChangeServlet(RestServlet):
    """Bloque le changement de mot de passe"""
    PATTERNS = [
        "^/_matrix/client/(r0|v3|unstable)/account/password$",
        "^/_matrix/client/(r0|v3|unstable)/account/deactivate$"
    ]

    def __init__(self, hs):
        super().__init__()
        self.hs = hs

    async def on_POST(self, request):
        raise SynapseError(
            403,
            "Changement de mot de passe et suppression de compte désactivés par l'administrateur",
            Codes.FORBIDDEN
        )

def parse_config(config):
    return config

def register_servlets(hs, http_server):
    BlockPasswordChangeServlet(hs).register(http_server)
