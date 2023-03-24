from datetime import datetime
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy.sql import func, text, expression
from sqlalchemy import (
    DateTime,
    BigInteger,
    Integer,
    Column,
    String,
    Float,
    Boolean,
    ForeignKey
)

Base = declarative_base()

class Product(Base):
    __tablename__ = 'product'

    product_uuid = Column(UUID(as_uuid=True), primary_key=True)
    name = Column(String, nullable=False)
    description = Column(String, nullable=False)
    price = Column(BigInteger, nullable=True)